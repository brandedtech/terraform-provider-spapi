package provider

func RegionToEndpoint(region string) string {
	switch region {
	case "NA":
		return "https://sellingpartnerapi-na.amazon.com"
	case "EU":
		return "https://sellingpartnerapi-eu.amazon.com"
	case "FE":
		return "https://sellingpartnerapi-fe.amazon.com"
	default:
		return "https://sellingpartnerapi-na.amazon.com"
	}
}
