package events;

import ( "testing" )

func  TestGetIDs(t *testing.T) {
	url := "https://discord.com/channels/123456789012345678/123456789012345678/123456789012345678"
	expected := []string{"123456789012345678", "123456789012345678", "123456789012345678"}
	actual := getIDs(url);
	if actual[0] != expected[0] || actual[1] != expected[1] || actual[2] != expected[2] {
		t.Error("Expected", expected, "got", actual)
	}
}
