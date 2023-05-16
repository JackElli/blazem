package cors

func GetAllowedMethods() []string {
	return []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
}

func GetAllowedHeaders() []string {
	return []string{"Origin", "Content-Length", "Content-Type"}
}

func GetAllowedOrigin() []string {
	return []string{"http://localhost:5173"}
}

func GetExposedHeaders() []string {
	return []string{"Set-Cookie"}
}
