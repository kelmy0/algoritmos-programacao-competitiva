import { SITEMAP_SECRET } from "$env/static/private";
import { PUBLIC_API_URL } from "$env/static/public";
import { customFetch } from "$lib/api/client";
import type { RequestHandler } from "./$types";

interface SitemapItem {
	slug: string;
	updatedAt: string;
}

interface ApiResponse {
	data: SitemapItem[];
}

export const GET: RequestHandler = async (event) => {
	const siteUrl = "https://algoritmos-programacao.com.br";

	const staticPages = [{ url: "", priority: "1.0", changefreq: "daily" }];

	let dynamicItems: SitemapItem[] = [];

	const { data, error, status } = await customFetch<ApiResponse>(
		event.fetch,
		`${PUBLIC_API_URL}/api/sitemap/algorithms`,
		{
			headers: {
				"x-sitemap-secret": SITEMAP_SECRET
			}
		}
	);

	if (error) {
		console.error("Error in Go API response when fetching sitemap:", status, error);
	} else if (data && data.data) {
		dynamicItems = data.data;
	}

	const body = `<?xml version="1.0" encoding="UTF-8" ?>
<urlset xmlns="https://www.sitemap.org/schemas/sitemap/0.9">
  <!-- Static Pages -->
  ${staticPages
		.map(
			(page) => `
  <url>
    <loc>${siteUrl}${page.url ? `/${page.url}` : ""}</loc>
    <changefreq>${page.changefreq}</changefreq>
    <priority>${page.priority}</priority>
  </url>`
		)
		.join("")}
  <!-- Dynamic Algorithms -->
  ${dynamicItems
		.map((item) => {
			const dateFormatted = item.updatedAt ? item.updatedAt.split("T")[0] : "";
			return `
  <url>
    <loc>${siteUrl}/algorithms/${item.slug}</loc>
    ${dateFormatted ? `<lastmod>${dateFormatted}</lastmod>` : ""}
    <changefreq>weekly</changefreq>
    <priority>0.8</priority>
  </url>`;
		})
		.join("")}
</urlset>`.trim();

	return new Response(body, {
		headers: {
			"Content-Type": "application/xml",
			"Cache-Control": "max-age=0, s-maxage=3600"
		}
	});
};
