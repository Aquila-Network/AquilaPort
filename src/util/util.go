package util

func getVersion(document Document) []byte {
	// version: timestamp (milliseconds, 13 digits) + deleted
	var delStatus byte
	delStatus = 48
	if document.Deleted {
		delStatus = 49
	}
	versionGen := append([]byte(document.Timestamp), delStatus)

	return versionGen
}
