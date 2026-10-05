package main

import (
	"fmt"
	"math"
)

// 1. Define an Interface
// In Go, an interface specifies a method set. Any type that implements
// these methods automatically implements the interface (no "implements" keyword needed).
type Shape interface {
	Area() float64
	Perimeter() float64
}

// 2. Concrete Type: Rectangle
type Rectangle struct {
	Width, Height float64
}

// Rectangle implements Shape by defining Area() and Perimeter()
func (r Rectangle) Area() float64 {
	return r.Width * r.Height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.Width + r.Height)
}

// 3. Concrete Type: Circle
type Circle struct {
	Radius float64
}

// Circle implements Shape
func (c Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.Radius
}

// 4. Polymorphism: A function that accepts any Shape
func printShapeInfo(s Shape) {
	fmt.Printf("Type: %T\n", s)
	fmt.Printf("Area: %.2f\n", s.Area())
	fmt.Printf("Perimeter: %.2f\n\n", s.Perimeter())
}

// 5. Type Assertions & Type Switches
// The empty interface `any` (or `interface{}`) can hold values of any type.
func describe(val any) {
	switch v := val.(type) {
	case int:
		fmt.Printf("It's an integer: %d\n", v)
	case string:
		fmt.Printf("It's a string: %q\n", v)
	case Shape:
		fmt.Printf("It's a Shape with Area: %.2f\n", v.Area())
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

func main() {
	fmt.Println("--- 1. Polymorphism with Interfaces ---")
	rect := Rectangle{Width: 10, Height: 5}
	circ := Circle{Radius: 7}

	// Both Rectangle and Circle can be passed to printShapeInfo
	printShapeInfo(rect)
	printShapeInfo(circ)

	// A slice of interfaces
	shapes := []Shape{rect, circ}
	var totalArea float64
	for _, shape := range shapes {
		totalArea += shape.Area()
	}
	fmt.Printf("Total Area of all shapes: %.2f\n\n", totalArea)

	fmt.Println("--- 2. Type Switch & Empty Interface (`any`) ---")
	describe(42)
	describe("Hello, Go!")
	describe(rect)
	describe(true)

	fmt.Println("\n--- 3. Type Assertion ---")
	var s Shape = Rectangle{Width: 4, Height: 6}

	// Extract concrete type with type assertion: value, ok := interfaceVar.(ConcreteType)
	r, ok := s.(Rectangle)
	if ok {
		fmt.Printf("Successfully asserted Rectangle: Width=%.1f, Height=%.1f\n", r.Width, r.Height)
	} else {
		fmt.Println("Failed to assert as Rectangle")
	}
}
