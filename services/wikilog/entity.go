package wikilog

import (
	elastic "gopkg.in/olivere/elastic.v5"
	"context"
	"reflect"
)

// Domain provides operations on wikiLog.
type Domain interface {
	searchPageName(string) ([]string, error)
}

type wikilogDomain struct{}

type wikiViews struct {
	Text     string
	Page     string
	Language string
	Sum      int
}

func (wikilogDomain) searchPageName(s string) ([]string, error) {
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		panic(err)
	}

	query := elastic.NewMatchQuery("text", s)

	searchResult, err := client.Search().
		Index("wikiviews").
		Query(query).
		From(0).Size(10).
		Pretty(true).
		Do(ctx)
	if err != nil {
		panic(err)
	}

	var result = make([]string, len(searchResult.Hits.Hits))

	var ttyp wikiViews
	for k, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(wikiViews); ok {
			result[k] = t.Text
		}
	}

	return result, err
}

// ServiceMiddleware is a chainable behavior modifier for WikilogDomain.
type ServiceMiddleware func(Domain) Domain
