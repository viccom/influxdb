package spectests

import (
	"time"

	"github.com/influxdata/flux"
	"github.com/influxdata/flux/ast"
	"github.com/influxdata/flux/execute"
	"github.com/influxdata/flux/semantic"
	"github.com/influxdata/flux/stdlib/influxdata/influxdb"
	"github.com/influxdata/flux/stdlib/universe"
	"github.com/influxdata/influxql"
)

func init() {
	RegisterFixture(
		NewFixture(
			`SELECT mean(value) FROM db0..cpu; SELECT max(value) FROM db0..cpu`,
			&flux.Spec{
				Operations: []*flux.Operation{
					{
						ID: "from0",
						Spec: &influxdb.FromOpSpec{
							BucketID: bucketID.String(),
						},
					},
					{
						ID: "range0",
						Spec: &universe.RangeOpSpec{
							Start:       flux.Time{Absolute: time.Unix(0, influxql.MinTime)},
							Stop:        flux.Time{Absolute: time.Unix(0, influxql.MaxTime)},
							TimeColumn:  execute.DefaultTimeColLabel,
							StartColumn: execute.DefaultStartColLabel,
							StopColumn:  execute.DefaultStopColLabel,
						},
					},
					{
						ID: "filter0",
						Spec: &universe.FilterOpSpec{
							Fn: &semantic.FunctionExpression{
								Block: &semantic.FunctionBlock{
									Parameters: &semantic.FunctionParameters{
										List: []*semantic.FunctionParameter{
											{Key: &semantic.Identifier{Name: "r"}},
										},
									},
									Body: &semantic.LogicalExpression{
										Operator: ast.AndOperator,
										Left: &semantic.BinaryExpression{
											Operator: ast.EqualOperator,
											Left: &semantic.MemberExpression{
												Object: &semantic.IdentifierExpression{
													Name: "r",
												},
												Property: "_measurement",
											},
											Right: &semantic.StringLiteral{
												Value: "cpu",
											},
										},
										Right: &semantic.BinaryExpression{
											Operator: ast.EqualOperator,
											Left: &semantic.MemberExpression{
												Object: &semantic.IdentifierExpression{
													Name: "r",
												},
												Property: "_field",
											},
											Right: &semantic.StringLiteral{
												Value: "value",
											},
										},
									},
								},
							},
						},
					},
					{
						ID: "group0",
						Spec: &universe.GroupOpSpec{
							Columns: []string{"_measurement", "_start"},
							Mode:    "by",
						},
					},
					{
						ID: "mean0",
						Spec: &universe.MeanOpSpec{
							AggregateConfig: execute.AggregateConfig{
								Columns: []string{execute.DefaultValueColLabel},
							},
						},
					},
					{
						ID: "duplicate0",
						Spec: &universe.DuplicateOpSpec{
							Column: execute.DefaultStartColLabel,
							As:     execute.DefaultTimeColLabel,
						},
					},
					{
						ID: "map0",
						Spec: &universe.MapOpSpec{
							Fn: &semantic.FunctionExpression{
								Block: &semantic.FunctionBlock{
									Parameters: &semantic.FunctionParameters{
										List: []*semantic.FunctionParameter{{
											Key: &semantic.Identifier{Name: "r"},
										}},
									},
									Body: &semantic.ObjectExpression{
										Properties: []*semantic.Property{
											{
												Key: &semantic.Identifier{Name: "_time"},
												Value: &semantic.MemberExpression{
													Object: &semantic.IdentifierExpression{
														Name: "r",
													},
													Property: "_time",
												},
											},
											{
												Key: &semantic.Identifier{Name: "mean"},
												Value: &semantic.MemberExpression{
													Object: &semantic.IdentifierExpression{
														Name: "r",
													},
													Property: "_value",
												},
											},
										},
									},
								},
							},
							MergeKey: true,
						},
					},
					{
						ID: "yield0",
						Spec: &universe.YieldOpSpec{
							Name: "0",
						},
					},
					{
						ID: "from1",
						Spec: &influxdb.FromOpSpec{
							BucketID: bucketID.String(),
						},
					},
					{
						ID: "range1",
						Spec: &universe.RangeOpSpec{
							Start:       flux.Time{Absolute: time.Unix(0, influxql.MinTime)},
							Stop:        flux.Time{Absolute: time.Unix(0, influxql.MaxTime)},
							TimeColumn:  execute.DefaultTimeColLabel,
							StartColumn: execute.DefaultStartColLabel,
							StopColumn:  execute.DefaultStopColLabel,
						},
					},
					{
						ID: "filter1",
						Spec: &universe.FilterOpSpec{
							Fn: &semantic.FunctionExpression{
								Block: &semantic.FunctionBlock{
									Parameters: &semantic.FunctionParameters{
										List: []*semantic.FunctionParameter{
											{Key: &semantic.Identifier{Name: "r"}},
										},
									},
									Body: &semantic.LogicalExpression{
										Operator: ast.AndOperator,
										Left: &semantic.BinaryExpression{
											Operator: ast.EqualOperator,
											Left: &semantic.MemberExpression{
												Object: &semantic.IdentifierExpression{
													Name: "r",
												},
												Property: "_measurement",
											},
											Right: &semantic.StringLiteral{
												Value: "cpu",
											},
										},
										Right: &semantic.BinaryExpression{
											Operator: ast.EqualOperator,
											Left: &semantic.MemberExpression{
												Object: &semantic.IdentifierExpression{
													Name: "r",
												},
												Property: "_field",
											},
											Right: &semantic.StringLiteral{
												Value: "value",
											},
										},
									},
								},
							},
						},
					},
					{
						ID: "group1",
						Spec: &universe.GroupOpSpec{
							Columns: []string{"_measurement", "_start"},
							Mode:    "by",
						},
					},
					{
						ID: "max0",
						Spec: &universe.MaxOpSpec{
							SelectorConfig: execute.SelectorConfig{
								Column: execute.DefaultValueColLabel,
							},
						},
					},
					{
						ID: "map1",
						Spec: &universe.MapOpSpec{
							Fn: &semantic.FunctionExpression{
								Block: &semantic.FunctionBlock{
									Parameters: &semantic.FunctionParameters{
										List: []*semantic.FunctionParameter{{
											Key: &semantic.Identifier{Name: "r"},
										}},
									},
									Body: &semantic.ObjectExpression{
										Properties: []*semantic.Property{
											{
												Key: &semantic.Identifier{Name: "_time"},
												Value: &semantic.MemberExpression{
													Object: &semantic.IdentifierExpression{
														Name: "r",
													},
													Property: "_time",
												},
											},
											{
												Key: &semantic.Identifier{Name: "max"},
												Value: &semantic.MemberExpression{
													Object: &semantic.IdentifierExpression{
														Name: "r",
													},
													Property: "_value",
												},
											},
										},
									},
								},
							},
							MergeKey: true,
						},
					},
					{
						ID: "yield1",
						Spec: &universe.YieldOpSpec{
							Name: "1",
						},
					},
				},
				Edges: []flux.Edge{
					{Parent: "from0", Child: "range0"},
					{Parent: "range0", Child: "filter0"},
					{Parent: "filter0", Child: "group0"},
					{Parent: "group0", Child: "mean0"},
					{Parent: "mean0", Child: "duplicate0"},
					{Parent: "duplicate0", Child: "map0"},
					{Parent: "map0", Child: "yield0"},
					{Parent: "from1", Child: "range1"},
					{Parent: "range1", Child: "filter1"},
					{Parent: "filter1", Child: "group1"},
					{Parent: "group1", Child: "max0"},
					{Parent: "max0", Child: "map1"},
					{Parent: "map1", Child: "yield1"},
				},
				Now: Now(),
			},
		),
	)
}
