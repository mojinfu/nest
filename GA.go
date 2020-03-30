package nest

import (
	"log"
	"math/rand"
)

func (this *GAStruct) logIfDebug(v ...interface{}) {
	if this.config.IfDebug {
		log.Println(v...)
	}
}

type GAStruct struct {
	binBounds  *BoundStruct
	config     *ConfigStruct
	population []*populationStruct
	randomSeed *rand.Rand
}
type placementsStruct struct {
	Placements [][]*PositionStruct
	rotation   []int
	fitness    float64
	paths      []*PolygonStruct
	area       float64
}

const defaultfitness float64 = 10000000

type populationStruct struct {
	placement []*PolygonStruct
	rotation  []int
	fitness   float64
	paths     []Path
	area      float64
}
type populationSlice []*populationStruct

func (this populationSlice) Len() int {
	return len(this)
}
func (this populationSlice) Less(i, j int) bool {
	return this[i].fitness < this[j].fitness
}
func (this populationSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

func (this *GAStruct) randomAngle(part *PolygonStruct) int {
	if part.isWart {
		return 0
	}
	var angleList []int = []int{}
	if this.config.Rotations < 1 {
		this.config.Rotations = 1
	}
	if len(part.AngleList) == 0 {
		for i := 0; i < this.config.Rotations; i++ {
			angleList = append(angleList, i*(360/this.config.Rotations))
		}
	} else {
		for index := range part.AngleList {
			angleList = append(angleList, int(part.AngleList[index]))
		}
	}

	shuffleArray := func(array []int) []int {
		for i := len(array) - 1; i > 0; i-- {
			var j = this.randomSeed.Intn(100) * (i + 1) / 100
			var temp = array[i]
			array[i] = array[j]
			array[j] = temp
		}
		return array
	}

	angleList = shuffleArray(angleList)

	for i := 0; i < len(angleList); i++ {
		_, bound := rotatePolygonA(part, angleList[i])

		// don't use obviously bad angles where the part doesn't fit in the bin
		if bound.width < this.binBounds.width && bound.height < this.binBounds.height {
			this.logIfDebug("angleList[i]:", angleList[i])
			return angleList[i]

		}
	}
	return 0
}
func (this *SVG) NewGeneticAlgorithm(adam []*PolygonStruct, bin *BinPolygonStruct, config *ConfigStruct) *GAStruct {
	myGA := &GAStruct{}
	//myGA.randomSeed = rand.New(rand.NewSource(time.Now().UnixNano()))//saya

	myGA.randomSeed = rand.New(rand.NewSource(3))
	myGA.config = config
	myGA.binBounds = getPolygonBounds(bin.myPolygon)
	//myGA.IfDebug = this.config.IfDebug
	//log.Println(this.Config.IfDebug)
	// population is an array of individuals. Each individual is a object representing the order of insertion and the angle each part is rotated
	var angles []int = []int{}
	for i := 0; i < len(adam); i++ {
		if !adam[i].isWart {
			angles = append(angles, myGA.randomAngle(adam[i]))
		} else {
			angles = append(angles, 0)
		}
	}

	myGA.population = []*populationStruct{&populationStruct{placement: adam, rotation: angles, fitness: defaultfitness}}

	for len(myGA.population) < config.PopulationSize {
		mutant := myGA.mutate(myGA.population[0])
		myGA.population = append(myGA.population, mutant)
	}
	return myGA
}
func (this *GAStruct) mutate(individual *populationStruct) *populationStruct {
	clone := &populationStruct{} //copy
	for index := range individual.placement {
		clone.placement = append(clone.placement, individual.placement[index])
	}
	for index := range individual.rotation {
		clone.rotation = append(clone.rotation, individual.rotation[index])
	}

	for i := 0; i < len(clone.placement); i++ {

		myRand := this.randomSeed.Intn(100)

		if myRand < this.config.MutationRate && !clone.placement[i].isWart {
			// swap current part with next part
			j := i + 1
			if j < len(clone.placement) {
				if !clone.placement[j].isWart {
					//此处是否可以开启？
					clone.placement[i], clone.placement[j] = clone.placement[j], clone.placement[i]
					clone.rotation[i], clone.rotation[j] = clone.rotation[j], clone.rotation[i]
				}
			}
		}
		myRand = this.randomSeed.Intn(100)
		if myRand < this.config.MutationRate && !clone.placement[i].isWart {
			if !clone.placement[i].isWart {
				clone.rotation[i] = this.randomAngle(clone.placement[i])
			} else {
				panic("")
			}
		}
	}
	return clone
}
