package influxql

import (
	"bytes"
	"container/heap"
	"fmt"
	"math"
	"sort"
)

/*
This file contains iterator implementations for each function call available
in InfluxQL. Call iterators are separated into two groups:

1. Map/reduce-style iterators - these are passed to IteratorCreator so that
   processing can be at the low-level storage and aggregates are returned.

2. Raw aggregate iterators - these require the full set of data for a window.
   These are handled by the select() function and raw points are streamed in
   from the low-level storage.

There are helpers to aid in building aggregate iterators. For simple map/reduce
iterators, you can use the reduceIterator types and pass a reduce function. This
reduce function is passed a previous and current value and the new timestamp,
value, and auxilary fields are returned from it.

For raw aggregate iterators, you can use the reduceSliceIterators which pass
in a slice of all points to the function and return a point. For more complex
iterator types, you may need to create your own iterators by hand.

Once your iterator is complete, you'll need to add it to the NewCallIterator()
function if it is to be available to IteratorCreators and add it to the select()
function to allow it to be included during planning.
*/

// NewCallIterator returns a new iterator for a Call.
func NewCallIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	name := opt.Expr.(*Call).Name
	switch name {
	case "count":
		return newCountIterator(input, opt)
	case "min":
		return newMinIterator(input, opt)
	case "max":
		return newMaxIterator(input, opt)
	case "sum":
		return newSumIterator(input, opt)
	case "first":
		return newFirstIterator(input, opt)
	case "last":
		return newLastIterator(input, opt)
	case "mean":
		return newMeanIterator(input, opt)
	default:
		return nil, fmt.Errorf("unsupported function call: %s", name)
	}
}

// newCountIterator returns an iterator for operating on a count() call.
func newCountIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	// FIXME: Wrap iterator in int-type iterator and always output int value.

	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatCountReduce}, nil
	case IntegerIterator:
		return &integerReduceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerCountReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported count iterator type: %T", input)
	}
}

// floatCountReduce returns the count of points.
func FloatCountReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil {
		return opt.StartTime, 1, nil
	}
	return prev.Time, prev.Value + 1, nil
}

// integerCountReduce returns the count of points.
func IntegerCountReduce(prev, curr *IntegerPoint, opt *ReduceOptions) (int64, int64, []interface{}) {
	if prev == nil {
		return opt.StartTime, 1, nil
	}
	return prev.Time, prev.Value + 1, nil
}

// newMinIterator returns an iterator for operating on a min() call.
func newMinIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatMinReduce}, nil
	case IntegerIterator:
		return &integerReduceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerMinReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported min iterator type: %T", input)
	}
}

// floatMinReduce returns the minimum value between prev & curr.
func FloatMinReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil || curr.Value < prev.Value || (curr.Value == prev.Value && curr.Time < prev.Time) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// integerMinReduce returns the minimum value between prev & curr.
func IntegerMinReduce(prev, curr *IntegerPoint, opt *ReduceOptions) (int64, int64, []interface{}) {
	if prev == nil || curr.Value < prev.Value || (curr.Value == prev.Value && curr.Time < prev.Time) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// newMaxIterator returns an iterator for operating on a max() call.
func newMaxIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatMaxReduce}, nil
	case IntegerIterator:
		return &integerReduceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerMaxReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported max iterator type: %T", input)
	}
}

// floatMaxReduce returns the maximum value between prev & curr.
func FloatMaxReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil || curr.Value > prev.Value || (curr.Value == prev.Value && curr.Time < prev.Time) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// integerMaxReduce returns the maximum value between prev & curr.
func IntegerMaxReduce(prev, curr *IntegerPoint, opt *ReduceOptions) (int64, int64, []interface{}) {
	if prev == nil || curr.Value > prev.Value || (curr.Value == prev.Value && curr.Time < prev.Time) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// newSumIterator returns an iterator for operating on a sum() call.
func newSumIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatSumReduce}, nil
	case IntegerIterator:
		return &integerReduceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerSumReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported sum iterator type: %T", input)
	}
}

// floatSumReduce returns the sum prev value & curr value.
func FloatSumReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil {
		return curr.Time, curr.Value, nil
	}
	return prev.Time, prev.Value + curr.Value, nil
}

// integerSumReduce returns the sum prev value & curr value.
func IntegerSumReduce(prev, curr *IntegerPoint, opt *ReduceOptions) (int64, int64, []interface{}) {
	if prev == nil {
		return curr.Time, curr.Value, nil
	}
	return prev.Time, prev.Value + curr.Value, nil
}

// newFirstIterator returns an iterator for operating on a first() call.
func newFirstIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatFirstReduce}, nil
	case IntegerIterator:
		return &integerReduceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerFirstReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported first iterator type: %T", input)
	}
}

// floatFirstReduce returns the first point sorted by time.
func FloatFirstReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil || curr.Time < prev.Time || (curr.Time == prev.Time && curr.Value > prev.Value) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// integerFirstReduce returns the first point sorted by time.
func IntegerFirstReduce(prev, curr *IntegerPoint, opt *ReduceOptions) (int64, int64, []interface{}) {
	if prev == nil || curr.Time < prev.Time || (curr.Time == prev.Time && curr.Value > prev.Value) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// newLastIterator returns an iterator for operating on a last() call.
func newLastIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatLastReduce}, nil
	case IntegerIterator:
		return &integerReduceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerLastReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported last iterator type: %T", input)
	}
}

// floatLastReduce returns the last point sorted by time.
func FloatLastReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil || curr.Time > prev.Time || (curr.Time == prev.Time && curr.Value > prev.Value) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// integerLastReduce returns the last point sorted by time.
func IntegerLastReduce(prev, curr *IntegerPoint, opt *ReduceOptions) (int64, int64, []interface{}) {
	if prev == nil || curr.Time > prev.Time || (curr.Time == prev.Time && curr.Value > prev.Value) {
		return curr.Time, curr.Value, curr.Aux
	}
	return prev.Time, prev.Value, prev.Aux
}

// NewDistinctIterator returns an iterator for operating on a distinct() call.
func NewDistinctIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatDistinctReduceSlice}, nil
	case IntegerIterator:
		return &integerReduceSliceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerDistinctReduceSlice}, nil
	case StringIterator:
		return &stringReduceSliceIterator{input: newBufStringIterator(input), opt: opt, fn: StringDistinctReduceSlice}, nil
	default:
		return nil, fmt.Errorf("unsupported distinct iterator type: %T", input)
	}
}

// floatDistinctReduceSlice returns the distinct value within a window.
func FloatDistinctReduceSlice(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
	m := make(map[float64]FloatPoint)
	for _, p := range a {
		if _, ok := m[p.Value]; !ok {
			m[p.Value] = p
		}
	}

	points := make([]FloatPoint, 0, len(m))
	for _, p := range m {
		points = append(points, FloatPoint{Time: p.Time, Value: p.Value})
	}
	sort.Sort(floatPoints(points))
	return points
}

// integerDistinctReduceSlice returns the distinct value within a window.
func IntegerDistinctReduceSlice(a []IntegerPoint, opt *ReduceOptions) []IntegerPoint {
	m := make(map[int64]IntegerPoint)
	for _, p := range a {
		if _, ok := m[p.Value]; !ok {
			m[p.Value] = p
		}
	}

	points := make([]IntegerPoint, 0, len(m))
	for _, p := range m {
		points = append(points, IntegerPoint{Time: p.Time, Value: p.Value})
	}
	sort.Sort(integerPoints(points))
	return points
}

// stringDistinctReduceSlice returns the distinct value within a window.
func StringDistinctReduceSlice(a []StringPoint, opt *ReduceOptions) []StringPoint {
	m := make(map[string]StringPoint)
	for _, p := range a {
		if _, ok := m[p.Value]; !ok {
			m[p.Value] = p
		}
	}

	points := make([]StringPoint, 0, len(m))
	for _, p := range m {
		points = append(points, StringPoint{Time: p.Time, Value: p.Value})
	}
	sort.Sort(stringPoints(points))
	return points
}

// newMeanIterator returns an iterator for operating on a mean() call.
func newMeanIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatMeanReduce}, nil
	case IntegerIterator:
		return &integerReduceFloatIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerMeanReduce}, nil
	default:
		return nil, fmt.Errorf("unsupported mean iterator type: %T", input)
	}
}

// floatMeanReduce returns the mean value within a window.
func FloatMeanReduce(prev, curr *FloatPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil {
		return opt.StartTime, curr.Value, nil
	}

	value := prev.Value * float64(prev.Aggregated)
	if curr.Aggregated > 1 {
		value += curr.Value * float64(curr.Aggregated)
		value /= float64(prev.Aggregated + curr.Aggregated)
	} else {
		value += curr.Value
		value /= float64(prev.Aggregated + 1)
	}
	return prev.Time, value, prev.Aux
}

// integerMeanReduceSlice returns the mean value within a window.
func IntegerMeanReduce(prev *FloatPoint, curr *IntegerPoint, opt *ReduceOptions) (int64, float64, []interface{}) {
	if prev == nil {
		return opt.StartTime, float64(curr.Value), nil
	}

	value := prev.Value * float64(prev.Aggregated)
	if curr.Aggregated > 1 {
		value += float64(curr.Value) * float64(curr.Aggregated)
		value /= float64(prev.Aggregated + curr.Aggregated)
	} else {
		value += float64(curr.Value)
		value /= float64(prev.Aggregated + 1)
	}
	return prev.Time, value, prev.Aux
}

// newMedianIterator returns an iterator for operating on a median() call.
func newMedianIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatMedianReduceSlice}, nil
	case IntegerIterator:
		return &integerReduceSliceFloatIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerMedianReduceSlice}, nil
	default:
		return nil, fmt.Errorf("unsupported median iterator type: %T", input)
	}
}

// floatMedianReduceSlice returns the median value within a window.
func FloatMedianReduceSlice(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
	if len(a) == 1 {
		return []FloatPoint{{Time: opt.StartTime, Value: a[0].Value}}
	}

	// OPTIMIZE(benbjohnson): Use getSortedRange() from v0.9.5.1.

	// Return the middle value from the points.
	// If there are an even number of points then return the mean of the two middle points.
	sort.Sort(floatPointsByValue(a))
	if len(a)%2 == 0 {
		lo, hi := a[len(a)/2-1], a[(len(a)/2)]
		return []FloatPoint{{Time: opt.StartTime, Value: lo.Value + (hi.Value-lo.Value)/2}}
	}
	return []FloatPoint{{Time: opt.StartTime, Value: a[len(a)/2].Value}}
}

// integerMedianReduceSlice returns the median value within a window.
func IntegerMedianReduceSlice(a []IntegerPoint, opt *ReduceOptions) []FloatPoint {
	if len(a) == 1 {
		return []FloatPoint{{Time: opt.StartTime, Value: float64(a[0].Value)}}
	}

	// OPTIMIZE(benbjohnson): Use getSortedRange() from v0.9.5.1.

	// Return the middle value from the points.
	// If there are an even number of points then return the mean of the two middle points.
	sort.Sort(integerPointsByValue(a))
	if len(a)%2 == 0 {
		lo, hi := a[len(a)/2-1], a[(len(a)/2)]
		return []FloatPoint{{Time: opt.StartTime, Value: float64(lo.Value) + float64(hi.Value-lo.Value)/2}}
	}
	return []FloatPoint{{Time: opt.StartTime, Value: float64(a[len(a)/2].Value)}}
}

// newStddevIterator returns an iterator for operating on a stddev() call.
func newStddevIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatStddevReduceSlice}, nil
	case IntegerIterator:
		return &integerReduceSliceFloatIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerStddevReduceSlice}, nil
	case StringIterator:
		return &stringReduceSliceIterator{input: newBufStringIterator(input), opt: opt, fn: StringStddevReduceSlice}, nil
	default:
		return nil, fmt.Errorf("unsupported stddev iterator type: %T", input)
	}
}

// floatStddevReduceSlice returns the stddev value within a window.
func FloatStddevReduceSlice(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
	// If there is only one point then return 0.
	if len(a) < 2 {
		return []FloatPoint{{Time: opt.StartTime, Nil: true}}
	}

	// Calculate the mean.
	var mean float64
	var count int
	for _, p := range a {
		if math.IsNaN(p.Value) {
			continue
		}
		count++
		mean += (p.Value - mean) / float64(count)
	}

	// Calculate the variance.
	var variance float64
	for _, p := range a {
		if math.IsNaN(p.Value) {
			continue
		}
		variance += math.Pow(p.Value-mean, 2)
	}
	return []FloatPoint{{
		Time:  opt.StartTime,
		Value: math.Sqrt(variance / float64(count-1)),
	}}
}

// integerStddevReduceSlice returns the stddev value within a window.
func IntegerStddevReduceSlice(a []IntegerPoint, opt *ReduceOptions) []FloatPoint {
	// If there is only one point then return 0.
	if len(a) < 2 {
		return []FloatPoint{{Time: opt.StartTime, Nil: true}}
	}

	// Calculate the mean.
	var mean float64
	var count int
	for _, p := range a {
		count++
		mean += (float64(p.Value) - mean) / float64(count)
	}

	// Calculate the variance.
	var variance float64
	for _, p := range a {
		variance += math.Pow(float64(p.Value)-mean, 2)
	}
	return []FloatPoint{{
		Time:  opt.StartTime,
		Value: math.Sqrt(variance / float64(count-1)),
	}}
}

// stringStddevReduceSlice always returns "".
func StringStddevReduceSlice(a []StringPoint, opt *ReduceOptions) []StringPoint {
	return []StringPoint{{Time: opt.StartTime, Value: ""}}
}

// newSpreadIterator returns an iterator for operating on a spread() call.
func newSpreadIterator(input Iterator, opt IteratorOptions) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: FloatSpreadReduceSlice}, nil
	case IntegerIterator:
		return &integerReduceSliceIterator{input: newBufIntegerIterator(input), opt: opt, fn: IntegerSpreadReduceSlice}, nil
	default:
		return nil, fmt.Errorf("unsupported spread iterator type: %T", input)
	}
}

// floatSpreadReduceSlice returns the spread value within a window.
func FloatSpreadReduceSlice(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
	// Find min & max values.
	min, max := a[0].Value, a[0].Value
	for _, p := range a[1:] {
		min = math.Min(min, p.Value)
		max = math.Max(max, p.Value)
	}
	return []FloatPoint{{Time: opt.StartTime, Value: max - min}}
}

// integerSpreadReduceSlice returns the spread value within a window.
func IntegerSpreadReduceSlice(a []IntegerPoint, opt *ReduceOptions) []IntegerPoint {
	// Find min & max values.
	min, max := a[0].Value, a[0].Value
	for _, p := range a[1:] {
		if p.Value < min {
			min = p.Value
		}
		if p.Value > max {
			max = p.Value
		}
	}
	return []IntegerPoint{{Time: opt.StartTime, Value: max - min}}
}

// newTopIterator returns an iterator for operating on a top() call.
func newTopIterator(input Iterator, opt IteratorOptions, n *NumberLiteral, tags []int) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: NewFloatTopReduceSliceFunc(int(n.Val), tags, opt.Interval)}, nil
	case IntegerIterator:
		return &integerReduceSliceIterator{input: newBufIntegerIterator(input), opt: opt, fn: NewIntegerTopReduceSliceFunc(int(n.Val), tags, opt.Interval)}, nil
	default:
		return nil, fmt.Errorf("unsupported top iterator type: %T", input)
	}
}

// newFloatTopReduceSliceFunc returns the top values within a window.
func NewFloatTopReduceSliceFunc(n int, tags []int, interval Interval) FloatReduceSliceFunc {
	return func(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
		// Filter by tags if they exist.
		if tags != nil {
			a = filterFloatByUniqueTags(a, tags, func(cur, p *FloatPoint) bool {
				return p.Value > cur.Value || (p.Value == cur.Value && p.Time < cur.Time)
			})
		}

		// If we ask for more elements than exist, restrict n to be the length of the array.
		size := n
		if size > len(a) {
			size = len(a)
		}

		// Construct a heap preferring higher values and breaking ties
		// based on the earliest time for a point.
		h := floatPointsSortBy(a, func(a, b *FloatPoint) bool {
			if a.Value != b.Value {
				return a.Value > b.Value
			}
			return a.Time < b.Time
		})
		heap.Init(h)

		// Pop the first n elements and then sort by time.
		points := make([]FloatPoint, 0, size)
		for i := 0; i < size; i++ {
			p := heap.Pop(h).(FloatPoint)
			points = append(points, p)
		}

		// Either zero out all values or sort the points by time
		// depending on if a time interval was given or not.
		if !interval.IsZero() {
			for i := range points {
				points[i].Time = opt.StartTime
			}
		} else {
			sort.Stable(floatPoints(points))
		}
		return points
	}
}

// newIntegerTopReduceSliceFunc returns the top values within a window.
func NewIntegerTopReduceSliceFunc(n int, tags []int, interval Interval) IntegerReduceSliceFunc {
	return func(a []IntegerPoint, opt *ReduceOptions) []IntegerPoint {
		// Filter by tags if they exist.
		if tags != nil {
			a = filterIntegerByUniqueTags(a, tags, func(cur, p *IntegerPoint) bool {
				return p.Value > cur.Value || (p.Value == cur.Value && p.Time < cur.Time)
			})
		}

		// If we ask for more elements than exist, restrict n to be the length of the array.
		size := n
		if size > len(a) {
			size = len(a)
		}

		// Construct a heap preferring higher values and breaking ties
		// based on the earliest time for a point.
		h := integerPointsSortBy(a, func(a, b *IntegerPoint) bool {
			if a.Value != b.Value {
				return a.Value > b.Value
			}
			return a.Time < b.Time
		})
		heap.Init(h)

		// Pop the first n elements and then sort by time.
		points := make([]IntegerPoint, 0, size)
		for i := 0; i < size; i++ {
			p := heap.Pop(h).(IntegerPoint)
			points = append(points, p)
		}

		// Either zero out all values or sort the points by time
		// depending on if a time interval was given or not.
		if !interval.IsZero() {
			for i := range points {
				points[i].Time = opt.StartTime
			}
		} else {
			sort.Stable(integerPoints(points))
		}
		return points
	}
}

// newBottomIterator returns an iterator for operating on a bottom() call.
func newBottomIterator(input Iterator, opt IteratorOptions, n *NumberLiteral, tags []int) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: NewFloatBottomReduceSliceFunc(int(n.Val), tags, opt.Interval)}, nil
	case IntegerIterator:
		return &integerReduceSliceIterator{input: newBufIntegerIterator(input), opt: opt, fn: NewIntegerBottomReduceSliceFunc(int(n.Val), tags, opt.Interval)}, nil
	default:
		return nil, fmt.Errorf("unsupported bottom iterator type: %T", input)
	}
}

// newFloatBottomReduceSliceFunc returns the bottom values within a window.
func NewFloatBottomReduceSliceFunc(n int, tags []int, interval Interval) FloatReduceSliceFunc {
	return func(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
		// Filter by tags if they exist.
		if tags != nil {
			a = filterFloatByUniqueTags(a, tags, func(cur, p *FloatPoint) bool {
				return p.Value < cur.Value || (p.Value == cur.Value && p.Time < cur.Time)
			})
		}

		// If we ask for more elements than exist, restrict n to be the length of the array.
		size := n
		if size > len(a) {
			size = len(a)
		}

		// Construct a heap preferring lower values and breaking ties
		// based on the earliest time for a point.
		h := floatPointsSortBy(a, func(a, b *FloatPoint) bool {
			if a.Value != b.Value {
				return a.Value < b.Value
			}
			return a.Time < b.Time
		})
		heap.Init(h)

		// Pop the first n elements and then sort by time.
		points := make([]FloatPoint, 0, size)
		for i := 0; i < size; i++ {
			p := heap.Pop(h).(FloatPoint)
			points = append(points, p)
		}

		// Either zero out all values or sort the points by time
		// depending on if a time interval was given or not.
		if !interval.IsZero() {
			for i := range points {
				points[i].Time = opt.StartTime
			}
		} else {
			sort.Stable(floatPoints(points))
		}
		return points
	}
}

// newIntegerBottomReduceSliceFunc returns the bottom values within a window.
func NewIntegerBottomReduceSliceFunc(n int, tags []int, interval Interval) IntegerReduceSliceFunc {
	return func(a []IntegerPoint, opt *ReduceOptions) []IntegerPoint {
		// Filter by tags if they exist.
		if tags != nil {
			a = filterIntegerByUniqueTags(a, tags, func(cur, p *IntegerPoint) bool {
				return p.Value < cur.Value || (p.Value == cur.Value && p.Time < cur.Time)
			})
		}

		// If we ask for more elements than exist, restrict n to be the length of the array.
		size := n
		if size > len(a) {
			size = len(a)
		}

		// Construct a heap preferring lower values and breaking ties
		// based on the earliest time for a point.
		h := integerPointsSortBy(a, func(a, b *IntegerPoint) bool {
			if a.Value != b.Value {
				return a.Value < b.Value
			}
			return a.Time < b.Time
		})
		heap.Init(h)

		// Pop the first n elements and then sort by time.
		points := make([]IntegerPoint, 0, size)
		for i := 0; i < size; i++ {
			p := heap.Pop(h).(IntegerPoint)
			points = append(points, p)
		}

		// Either zero out all values or sort the points by time
		// depending on if a time interval was given or not.
		if !interval.IsZero() {
			for i := range points {
				points[i].Time = opt.StartTime
			}
		} else {
			sort.Stable(integerPoints(points))
		}
		return points
	}
}

func filterFloatByUniqueTags(a []FloatPoint, tags []int, cmpFunc func(cur, p *FloatPoint) bool) []FloatPoint {
	pointMap := make(map[string]FloatPoint)
	for _, p := range a {
		keyBuf := bytes.NewBuffer(nil)
		for i, index := range tags {
			if i > 0 {
				keyBuf.WriteString(",")
			}
			fmt.Fprintf(keyBuf, "%s", p.Aux[index])
		}
		key := keyBuf.String()

		cur, ok := pointMap[key]
		if ok {
			if cmpFunc(&cur, &p) {
				pointMap[key] = p
			}
		} else {
			pointMap[key] = p
		}
	}

	// Recreate the original array with our new filtered list.
	points := make([]FloatPoint, 0, len(pointMap))
	for _, p := range pointMap {
		points = append(points, p)
	}
	return points
}

func filterIntegerByUniqueTags(a []IntegerPoint, tags []int, cmpFunc func(cur, p *IntegerPoint) bool) []IntegerPoint {
	pointMap := make(map[string]IntegerPoint)
	for _, p := range a {
		keyBuf := bytes.NewBuffer(nil)
		for i, index := range tags {
			if i > 0 {
				keyBuf.WriteString(",")
			}
			fmt.Fprintf(keyBuf, "%s", p.Aux[index])
		}
		key := keyBuf.String()

		cur, ok := pointMap[key]
		if ok {
			if cmpFunc(&cur, &p) {
				pointMap[key] = p
			}
		} else {
			pointMap[key] = p
		}
	}

	// Recreate the original array with our new filtered list.
	points := make([]IntegerPoint, 0, len(pointMap))
	for _, p := range pointMap {
		points = append(points, p)
	}
	return points
}

// newPercentileIterator returns an iterator for operating on a percentile() call.
func newPercentileIterator(input Iterator, opt IteratorOptions, percentile float64) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: NewFloatPercentileReduceSliceFunc(percentile)}, nil
	case IntegerIterator:
		return &integerReduceSliceIterator{input: newBufIntegerIterator(input), opt: opt, fn: NewIntegerPercentileReduceSliceFunc(percentile)}, nil
	default:
		return nil, fmt.Errorf("unsupported percentile iterator type: %T", input)
	}
}

// newFloatPercentileReduceSliceFunc returns the percentile value within a window.
func NewFloatPercentileReduceSliceFunc(percentile float64) FloatReduceSliceFunc {
	return func(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
		length := len(a)
		i := int(math.Floor(float64(length)*percentile/100.0+0.5)) - 1

		if i < 0 || i >= length {
			return []FloatPoint{{Time: opt.StartTime, Nil: true}}
		}

		sort.Sort(floatPointsByValue(a))
		return []FloatPoint{{Time: opt.StartTime, Value: a[i].Value}}
	}
}

// newIntegerPercentileReduceSliceFunc returns the percentile value within a window.
func NewIntegerPercentileReduceSliceFunc(percentile float64) IntegerReduceSliceFunc {
	return func(a []IntegerPoint, opt *ReduceOptions) []IntegerPoint {
		length := len(a)
		i := int(math.Floor(float64(length)*percentile/100.0+0.5)) - 1

		if i < 0 || i >= length {
			return nil
		}

		sort.Sort(integerPointsByValue(a))
		return []IntegerPoint{{Time: opt.StartTime, Value: a[i].Value}}
	}
}

// newDerivativeIterator returns an iterator for operating on a derivative() call.
func newDerivativeIterator(input Iterator, opt IteratorOptions, interval Interval, isNonNegative bool) (Iterator, error) {
	switch input := input.(type) {
	case FloatIterator:
		return &floatReduceSliceIterator{input: newBufFloatIterator(input), opt: opt, fn: NewFloatDerivativeReduceSliceFunc(interval, isNonNegative)}, nil
	case IntegerIterator:
		return &integerReduceSliceFloatIterator{input: newBufIntegerIterator(input), opt: opt, fn: NewIntegerDerivativeReduceSliceFunc(interval, isNonNegative)}, nil
	default:
		return nil, fmt.Errorf("unsupported derivative iterator type: %T", input)
	}
}

// newFloatDerivativeReduceSliceFunc returns the derivative value within a window.
func NewFloatDerivativeReduceSliceFunc(interval Interval, isNonNegative bool) FloatReduceSliceFunc {
	prev := FloatPoint{Time: -1}

	return func(a []FloatPoint, opt *ReduceOptions) []FloatPoint {
		if len(a) == 0 {
			return a
		} else if len(a) == 1 {
			return []FloatPoint{{Time: a[0].Time, Nil: true}}
		}

		if prev.Time == -1 {
			prev = a[0]
		}

		output := make([]FloatPoint, 0, len(a)-1)
		for i := 1; i < len(a); i++ {
			p := &a[i]

			// Calculate the derivative of successive points by dividing the
			// difference of each value by the elapsed time normalized to the interval.
			diff := p.Value - prev.Value
			elapsed := p.Time - prev.Time

			value := 0.0
			if elapsed > 0 {
				value = diff / (float64(elapsed) / float64(interval.Duration))
			}

			prev = *p

			// Drop negative values for non-negative derivatives.
			if isNonNegative && diff < 0 {
				continue
			}

			output = append(output, FloatPoint{Time: p.Time, Value: value})
		}
		return output
	}
}

// newIntegerDerivativeReduceSliceFunc returns the derivative value within a window.
func NewIntegerDerivativeReduceSliceFunc(interval Interval, isNonNegative bool) IntegerReduceSliceFloatFunc {
	prev := IntegerPoint{Time: -1}

	return func(a []IntegerPoint, opt *ReduceOptions) []FloatPoint {
		if len(a) == 0 {
			return []FloatPoint{}
		} else if len(a) == 1 {
			return []FloatPoint{{Time: a[0].Time, Nil: true}}
		}

		if prev.Time == -1 {
			prev = a[0]
		}

		output := make([]FloatPoint, 0, len(a)-1)
		for i := 1; i < len(a); i++ {
			p := &a[i]

			// Calculate the derivative of successive points by dividing the
			// difference of each value by the elapsed time normalized to the interval.
			diff := float64(p.Value - prev.Value)
			elapsed := p.Time - prev.Time

			value := 0.0
			if elapsed > 0 {
				value = diff / (float64(elapsed) / float64(interval.Duration))
			}

			prev = *p

			// Drop negative values for non-negative derivatives.
			if isNonNegative && diff < 0 {
				continue
			}

			output = append(output, FloatPoint{Time: p.Time, Value: value})
		}
		return output
	}
}

type integerReduceFloatIterator struct {
	input  *bufIntegerIterator
	fn     integerReduceFloatFunc
	opt    IteratorOptions
	points []*FloatPoint
}

// Close closes the iterator and all child iterators.
func (itr *integerReduceFloatIterator) Close() error { return itr.input.Close() }

// Next returns the minimum value for the next available interval.
func (itr *integerReduceFloatIterator) Next() *FloatPoint {
	// Calculate next window if we have no more points.
	if len(itr.points) == 0 {
		itr.points = itr.reduce()
		if len(itr.points) == 0 {
			return nil
		}
	}

	// Pop next point off the stack.
	p := itr.points[len(itr.points)-1]
	itr.points = itr.points[:len(itr.points)-1]
	return p
}

// reduce executes fn once for every point in the next window.
// The previous value for the dimension is passed to fn.
func (itr *integerReduceFloatIterator) reduce() []*FloatPoint {
	// Calculate next window.
	startTime, endTime := itr.opt.Window(itr.input.peekTime())

	var reduceOptions = ReduceOptions{
		StartTime: startTime,
		EndTime:   endTime,
	}

	// Create points by tags.
	m := make(map[string]*FloatPoint)
	for {
		// Read next point.
		curr := itr.input.NextInWindow(startTime, endTime)
		if curr == nil {
			break
		} else if curr.Nil {
			continue
		}
		tags := curr.Tags.Subset(itr.opt.Dimensions)
		id := curr.Name + "\x00" + tags.ID()

		// Pass previous and current points to reducer.
		prev := m[id]
		t, v, aux := itr.fn(prev, curr, &reduceOptions)
		if t == ZeroTime {
			continue
		}

		// If previous value didn't exist, create it and copy values.
		if prev == nil {
			prev = &FloatPoint{Name: curr.Name, Tags: tags}
			m[id] = prev
		}
		prev.Time = t
		prev.Value = v
		prev.Aux = aux
		prev.Aggregated++
	}

	// Reverse sort points by name & tag.
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	a := make([]*FloatPoint, len(m))
	for i, k := range keys {
		a[i] = m[k]
	}

	// Set the time on each point to the beginning of the interval.
	for _, p := range a {
		p.Time = startTime
	}

	return a
}

type integerReduceFloatFunc func(prev *FloatPoint, curr *IntegerPoint, opt *ReduceOptions) (int64, float64, []interface{})

// integerReduceSliceFloatIterator executes a reducer on all points in a window and buffers the result.
// This iterator receives an integer iterator but produces a float iterator.
type integerReduceSliceFloatIterator struct {
	input  *bufIntegerIterator
	fn     IntegerReduceSliceFloatFunc
	opt    IteratorOptions
	points []FloatPoint
}

// Close closes the iterator and all child iterators.
func (itr *integerReduceSliceFloatIterator) Close() error { return itr.input.Close() }

// Next returns the minimum value for the next available interval.
func (itr *integerReduceSliceFloatIterator) Next() *FloatPoint {
	// Calculate next window if we have no more points.
	if len(itr.points) == 0 {
		itr.points = itr.reduce()
		if len(itr.points) == 0 {
			return nil
		}
	}

	// Pop next point off the stack.
	p := itr.points[len(itr.points)-1]
	itr.points = itr.points[:len(itr.points)-1]
	return &p
}

// reduce executes fn once for every point in the next window.
// The previous value for the dimension is passed to fn.
func (itr *integerReduceSliceFloatIterator) reduce() []FloatPoint {
	// Calculate next window.
	startTime, endTime := itr.opt.Window(itr.input.peekTime())

	var reduceOptions = ReduceOptions{
		StartTime: startTime,
		EndTime:   endTime,
	}

	// Group points by name and tagset.
	groups := make(map[string]struct {
		name   string
		tags   Tags
		points []IntegerPoint
	})
	for {
		// Read next point.
		p := itr.input.NextInWindow(startTime, endTime)
		if p == nil {
			break
		}
		tags := p.Tags.Subset(itr.opt.Dimensions)

		// Append point to dimension.
		id := tags.ID()
		g := groups[id]
		g.name = p.Name
		g.tags = tags
		g.points = append(g.points, *p)
		groups[id] = g
	}

	// Reduce each set into a set of values.
	results := make(map[string][]FloatPoint)
	for key, g := range groups {
		a := itr.fn(g.points, &reduceOptions)
		if len(a) == 0 {
			continue
		}

		// Update name and tags for each returned point.
		for i := range a {
			a[i].Name = g.name
			a[i].Tags = g.tags
		}
		results[key] = a
	}

	// Reverse sort points by name & tag.
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	// Reverse order points within each key.
	a := make([]FloatPoint, 0, len(results))
	for _, k := range keys {
		for i := len(results[k]) - 1; i >= 0; i-- {
			a = append(a, results[k][i])
		}
	}

	return a
}

// integerReduceSliceFloatFunc is the function called by a IntegerPoint slice reducer that emits FloatPoint.
type IntegerReduceSliceFloatFunc func(a []IntegerPoint, opt *ReduceOptions) []FloatPoint
