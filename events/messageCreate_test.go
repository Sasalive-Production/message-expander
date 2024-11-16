package events;

import ( "testing" )

func  TestGetIDs(t *testing.T) {
	url := "https://discord.com/channels/123456789012345678/123456789012345678/123456789012345678"
	expected := messageInfo{guild: "123456789012345678", channel: "123456789012345678", message: "123456789012345678"}
	actual := getIDs(url);
	if actual != expected {
		t.Error("Expected", expected, "got", actual)
	}
}
