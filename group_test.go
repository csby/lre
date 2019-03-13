package lre

import "testing"

func TestGroups_Add(t *testing.T) {
	groups := make(Groups)

	group := groups.Add("m1")
	if group.Model != "m1" {
		t.Fatalf("group model except = 'm1', actual = '%s'", group.Model)
	}
	if len(groups) != 1 {
		t.Fatal("group count except = 1, actual = ", len(groups))
	}

	group = groups.Add("m2")
	if group.Model != "m2" {
		t.Fatalf("group model except = 'm2', actual = '%s'", group.Model)
	}
	if len(groups) != 2 {
		t.Fatal("group count except = 2, actual = ", len(groups))
	}

	group = groups.Add("m3")
	if group.Model != "m3" {
		t.Fatalf("group model except = 'm3', actual = '%s'", group.Model)
	}
	if len(groups) != 3 {
		t.Fatal("group count except = 2, actual = ", len(groups))
	}

	group = groups.Add("m2")
	if group.Model != "m2" {
		t.Fatalf("group model except = 'm2', actual = '%s'", group.Model)
	}
	if len(groups) != 3 {
		t.Fatal("group count except = 3, actual = ", len(groups))
	}
}
