export interface DeviceHeaders {
	"User-Agent": string;
	"Accept-Language": string;
	"Sec-CH-UA-Platform": string;
	"Sec-CH-UA-Mobile": string;
}

export function extractDeviceHeaders(request: Request): DeviceHeaders {
	return {
		"User-Agent": request.headers.get("user-agent") || "",
		"Accept-Language": request.headers.get("accept-language") || "",
		"Sec-CH-UA-Platform": request.headers.get("sec-ch-ua-platform") || "",
		"Sec-CH-UA-Mobile": request.headers.get("sec-ch-ua-mobile") || ""
	};
}
