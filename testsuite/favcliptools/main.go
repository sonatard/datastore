//go:generate jwg -output model_json.go -transcripttag swagger .
//go:generate qbg -output model_query.go -usedatastorewrapper .

package favcliptools

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go.mercari.io/datastore"
	"go.mercari.io/datastore/boom"
	"go.mercari.io/datastore/testsuite"
)

var _ datastore.PropertyTranslator = userID(0)
var _ json.Marshaler = userID(1)
var _ json.Unmarshaler = (*userID)(nil)

type contextClient struct{}

const kindUser = "User"

type userID int64

// User kind
// +jwg
// +qbg
type User struct {
	ID       userID `datastore:"-" boom:"id" json:"id"`
	Name     string `json:"name"`
	MentorID userID `json:"mentorID"`
}

func (id userID) ToPropertyValue(ctx context.Context) (interface{}, error) {
	client := ctx.Value(contextClient{}).(datastore.Client)
	key := client.IDKey(kindUser, int64(id), nil)
	return key, nil
}

func (id userID) FromPropertyValue(ctx context.Context, p datastore.Property) (dst interface{}, err error) {
	key, ok := p.Value.(datastore.Key)
	if !ok {
		return nil, datastore.ErrInvalidEntityType
	}
	return userID(key.ID()), nil
}

func (id userID) MarshalJSON() ([]byte, error) {
	jsonNumber := json.Number(fmt.Sprintf("%d", int64(id)))
	return json.Marshal(jsonNumber)
}

func (id *userID) UnmarshalJSON(b []byte) error {
	var jsonNumber json.Number
	err := json.Unmarshal(b, &jsonNumber)
	if err != nil {
		return err
	}
	v, err := jsonNumber.Int64()
	if err != nil {
		return err
	}

	*id = userID(v)
	return nil
}

// TestSuite contains all the test cases that this package provides.
var TestSuite = map[string]testsuite.Test{
	"FavclipTools": favclipTools,
}

func init() {
	testsuite.MergeTestSuite(TestSuite)
}

func favclipTools(ctx context.Context, t *testing.T, client datastore.Client) {
	defer func() {
		err := client.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// qbg using client.Context internally
	ctx = context.WithValue(ctx, contextClient{}, client)
	client.SetContext(ctx)

	bm := boom.FromClient(ctx, client)

	user := &User{
		ID:       userID(100),
		Name:     "foobar",
		MentorID: userID(200),
	}

	_, err := bm.Put(user)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.Marshal(user)
	if err != nil {
		t.Fatal(err)
	}

	if v := string(b); v != `{"id":100,"name":"foobar","mentorID":200}` {
		t.Errorf("unexpected: %v", v)
	}

	user = &User{}
	err = json.Unmarshal(b, user)
	if err != nil {
		t.Fatal(err)
	}

	if v := int64(user.ID); v != 100 {
		t.Errorf("unexpected: %v", v)
	}
	if v := user.Name; v != "foobar" {
		t.Errorf("unexpected: %v", v)
	}
	if v := int64(user.MentorID); v != 200 {
		t.Errorf("unexpected: %v", v)
	}

	{ // for jwg
		b := NewUserJSONBuilder()
		b.Add(b.ID)
		b.Add(b.MentorID)

		userJSON, err := b.Convert(user)
		if err != nil {
			t.Fatal(err)
		}

		if v := int64(userJSON.ID); v != 100 {
			t.Errorf("unexpected: %v", v)
		}
		// removed
		if v := userJSON.Name; v != "" {
			t.Errorf("unexpected: %v", v)
		}
		if v := int64(userJSON.MentorID); v != 200 {
			t.Errorf("unexpected: %v", v)
		}
	}
	{ // for qbg
		b := NewUserQueryBuilder(client)
		b.MentorID.Equal(userID(200))
		var list []*User
		_, err = bm.GetAll(b.Query(), &list)
		if err != nil {
			t.Fatal(err)
		}

		if v := len(list); v != 1 {
			t.Errorf("unexpected: %v", v)
		}
	}
}
