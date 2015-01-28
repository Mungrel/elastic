// Copyright 2012-2014 Oliver Eilhard. All rights reserved.
// Use of this source code is governed by a MIT-license.
// See http://olivere.mit-license.org/license.txt for details.

package elastic

// A boosting query can be used to effectively
// demote results that match a given query.
// For more details, see:
// http://www.elasticsearch.org/guide/en/elasticsearch/reference/current/query-dsl-boosting-query.html
type BoostingQuery struct {
	Query
	positiveClause []Query
	negativeClause []Query
	negativeBoost   *float32
}

// Creates a new boosting query.
func NewBoostingQuery() BoostingQuery {
	q := BoostingQuery{
		positiveClause:  make([]Query, 0),
		negativeClause: make([]Query, 0),
	}
	return q
}

func (q BoostingQuery) Positive(positive ...Query) BoostingQuery {
	q.positiveClause = positive
	return q
}

func (q BoostingQuery) Negative(negative ...Query) BoostingQuery {
	q.negativeClause = negative
	return q
}

func (q BoostingQuery) NegativeBoost(negativeBoost float32) BoostingQuery {
	q.negativeBoost = &negativeBoost
	return q
}

// Creates the query source for the boosting query.
func (q BoostingQuery) Source() interface{} {
	// {
	//     "boosting" : {
	//         "positive" : {
	//             "term" : {
	//                 "field1" : "value1"
	//             }
	//         },
	//         "negative" : {
	//             "term" : {
	//                 "field2" : "value2"
	//             }
	//         },
	//         "negative_boost" : 0.2
	//     }
	// }

	query := make(map[string]interface{})

	boostingClause := make(map[string]interface{})
	query["boosting"] = boostingClause

	// positive
	if len(q.positiveClause) == 1 {
		boostingClause["positive"] = q.positiveClause[0].Source()
	} else if len(q.positiveClause) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.positiveClause {
			clauses = append(clauses, subQuery.Source())
		}
		boostingClause["positive"] = clauses
	}

	// negative
	if len(q.negativeClause) == 1 {
		boostingClause["negative"] = q.negativeClause[0].Source()
	} else if len(q.negativeClause) > 1 {
		clauses := make([]interface{}, 0)
		for _, subQuery := range q.negativeClause {
			clauses = append(clauses, subQuery.Source())
		}
		boostingClause["negative"] = clauses
	}

	if q.negativeBoost != nil {
		boostingClause["negative_boost"] = *q.negativeBoost
	}

	return query
}
