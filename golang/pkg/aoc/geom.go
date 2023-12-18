package aoc

import "image"

// IsPointInsideLoop returns true if the targetPoint is inside the loop.
// ref: https://en.wikipedia.org/wiki/Point_in_polygon#Ray_casting_algorithm
func IsPointInsideLoop(targetPoint image.Point, loop []image.Point, allowOnBorder bool) bool {
	intersectionCount := 0

	for i := 0; i < len(loop); i++ {
		currentPoint := loop[i]
		nextPoint := loop[(i+1)%len(loop)]

		if allowOnBorder {
			isOnVertical := currentPoint.X == nextPoint.X && targetPoint.X == currentPoint.X
			isOnHorizontal := currentPoint.Y == nextPoint.Y && targetPoint.Y == currentPoint.Y

			if isOnVertical {
				if (targetPoint.Y >= currentPoint.Y && targetPoint.Y <= nextPoint.Y) || (targetPoint.Y <= currentPoint.Y && targetPoint.Y >= nextPoint.Y) {
					return true
				}
			}
			if isOnHorizontal {
				if (targetPoint.X >= currentPoint.X && targetPoint.X <= nextPoint.X) || (targetPoint.X <= currentPoint.X && targetPoint.X >= nextPoint.X) {
					return true
				}
			}
		}

		isAboveCurrent := currentPoint.Y > targetPoint.Y
		isAboveNext := nextPoint.Y > targetPoint.Y

		// If one point is above the targetPoint and the other point is below the targetPoint.
		// Then the ray from the targetPoint intersects the edge between currentPoint and nextPoint.
		isTargetBelowIntersection := isAboveCurrent != isAboveNext

		if !isTargetBelowIntersection {
			continue
		}

		nextDeltaX := nextPoint.X - currentPoint.X
		nextDeltaY := nextPoint.Y - currentPoint.Y
		targetDeltaY := targetPoint.Y - currentPoint.Y

		// Find the X-coordinate where the ray from the targetPoint intersects the edge between currentPoint and nextPoint.
		// If the targetPoint is to the left of the intersection, the ray intersects the edge.
		xIntersection := (nextDeltaX * targetDeltaY / nextDeltaY) + currentPoint.X
		isTargetLeftOfIntersection := targetPoint.X < xIntersection

		if isTargetLeftOfIntersection {
			intersectionCount++
		}
	}

	// If the number of intersections is odd, the targetPoint is inside the loop.
	return intersectionCount%2 == 1
}

func CalcSimplePolygonArea(points []image.Point) float64 {
	area := 0.0
	numPoints := len(points)

	for i := 0; i < numPoints; i++ {
		curr := points[i]
		next := points[(i+1)%numPoints] // Get the next point, wrap around to the first point if it's the last one

		area += float64(curr.X*next.Y - curr.Y*next.X)
	}

	return 0.5 * area
}
