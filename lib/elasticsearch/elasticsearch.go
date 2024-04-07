package elasticsearch

import (
	"anhle-grpc/model"
	"anhle-grpc/util"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/elastic/go-elasticsearch/esapi"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
)

type (
	IElasticsearch interface {
		GetClient() *elasticsearch.Client
	}

	Config struct {
		Username              string
		Password              string
		Host                  []string
		RetryStatuses         []int
		MaxRetries            int
		ResponseHeaderTimeout int
	}

	elasticsearchClient struct {
		config Config
		client *elasticsearch.Client
	}

	Response struct {
		StatusCode int
		Header     http.Header
		Body       map[string]any
	}
)

func NewElasticsearch(config Config) IElasticsearch {
	es := &elasticsearchClient{
		config: config,
	}
	if err := es.Connect(); err != nil {
		log.Fatal(err)
		return nil
	}
	_, err := es.GetClient().Ping()
	if err != nil {
		log.Fatal("Elasticsearch connection failed")
		return nil
	} else {
		log.Info("Elasticsearch connection successful")
	}
	return es
}

func (e *elasticsearchClient) Connect() error {
	client, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses:     e.config.Host,
			Username:      e.config.Username,
			Password:      e.config.Password,
			RetryOnStatus: e.config.RetryStatuses,
			MaxRetries:    e.config.MaxRetries,
		},
	)
	if err != nil {
		return err
	}
	e.client = client
	return nil
}

func (e *elasticsearchClient) GetClient() *elasticsearch.Client {
	return e.client
}

func ParseAnyToAny(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func DecodeAny(value io.Reader, dest any) error {
	if err := json.NewDecoder(value).Decode(dest); err != nil {
		return err
	}
	return nil
}

func EncodeAny(value any) (bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(value); err != nil {
		return buf, err
	}
	return buf, nil
}

func RangeQuery(field string, from, to any) map[string]any {
	return map[string]any{
		"range": map[string]any{
			field: map[string]any{
				"from":          from,
				"to":            to,
				"include_lower": true,
				"include_upper": true,
			},
		},
	}
}
func MatchQuery(field string, value any) map[string]any {
	return map[string]any{
		"match": map[string]any{
			field: value,
		},
	}
}

func TermQuery(field string, value any) map[string]any {
	return map[string]any{
		"term": map[string]any{
			field: map[string]any{
				"value": value,
			},
		},
	}
}

func TermsQuery(field string, values ...any) map[string]any {
	arr := make([]any, 0)
	if len(values) > 0 {
		arr = append(arr, values...)
	}
	return map[string]any{
		"terms": map[string]any{
			field: arr,
		},
	}
}

func WildcardQuery(field string, value any, insensitive sql.NullBool) map[string]any {
	valueTmp := map[string]any{
		"value": value,
	}
	if insensitive.Valid {
		valueTmp["case_insensitive"] = insensitive.Bool
	}
	return map[string]any{
		"wildcard": map[string]any{
			field: valueTmp,
		},
	}
}

func ShouldQuery(queries ...map[string]any) map[string]any {
	return map[string]any{
		"should": queries,
	}
}

func BoolQuery(queries ...map[string]any) map[string]any {
	query := map[string]any{
		"bool": map[string]any{},
	}
	if len(queries) == 1 {
		query["bool"] = queries[0]
	} else if len(queries) > 1 {
		var clauses []any
		for _, subQuery := range queries {
			clauses = append(clauses, subQuery)
		}
		query["bool"] = clauses
	}
	return query
}

func MustQuery(queries ...map[string]any) map[string]any {
	query := map[string]any{
		"must": map[string]any{},
	}
	if len(queries) == 1 {
		query["must"] = queries[0]
	} else if len(queries) > 1 {
		var clauses []any
		for _, subQuery := range queries {
			clauses = append(clauses, subQuery)
		}
		query["must"] = clauses
	}
	return query
}

func MustNotQuery(queries ...map[string]any) map[string]any {
	query := map[string]any{
		"must_not": map[string]any{},
	}
	if len(queries) == 1 {
		query["must_not"] = queries[0]
	} else if len(queries) > 1 {
		var clauses []any
		for _, subQuery := range queries {
			clauses = append(clauses, subQuery)
		}
		query["must_not"] = clauses
	}
	return query
}

func Order(field string, isAsc bool) map[string]any {
	order := "asc"
	if !isAsc {
		order = "desc"
	}
	return map[string]any{
		field: map[string]any{
			"order": order,
		},
	}
}

func ParseSearchResponse(response *esapi.Response) (*model.SearchReponse, error) {
	if response.IsError() {
		var e map[string]any
		if err := json.NewDecoder(response.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			typeErr := e["error"].(map[string]any)["type"]
			reason := e["error"].(map[string]any)["reason"]
			return nil, errors.New(util.ParseString(typeErr) + ":" + util.ParseString(reason))
		}
	}
	result := model.SearchReponse{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return &result, err
	}
	return &result, nil
}
