import { test, expect } from "@playwright/test";
import { readFileSync, readdirSync } from "fs";
import path from "path";

test.skip("index page", async ({ page }) => {
  await page.goto("/");
  await page.waitForLoadState("networkidle");
  await expect(page).toHaveTitle(/Hello, I'm Mark/);

  const links = [
    { text: "posts", href: "/posts" },
    { text: "polyring", href: "/polyring" },
    { text: "who is mark?", href: "/about" },
  ];

  for (const link of links) {
    await page.goto("/");

    const linkElement = await page.waitForSelector(
      `a[href="${link.href}"]:has-text("${link.text}")`,
    );
    await linkElement.click();
    await page.waitForLoadState("networkidle");

    await expect(page).toHaveURL(link.href);
  }
});

const files = readdirSync("content/posts");
const posts: { name: string; href: string }[] = [];
for (const f of files) {
  if (path.extname(f) !== ".md") continue;
  const file = readFileSync("content/posts/" + f, "utf-8");
  const name = file.split("\n")[1].replace("title: ", "").replace(/"/g, "");
  const href = "/posts/" + f.replace(".md", "").replace(/"/g, "");
  posts.push({ name, href });
}

for (const post of posts) {
  test.skip(`all posts to post: ${post.name}`, async ({ page }) => {
    await page.goto("/posts");

    const linkElement = await page.waitForSelector(
      `a[href="${post.href}"]:has-text("${post.name}")`,
    );
    const linkText = await linkElement?.textContent();
    expect(linkText && linkText.length > 0);

    await linkElement.click();
    await page.waitForLoadState("networkidle");

    const actualTitle = await page.title();
    expect(post.name).toEqual(actualTitle);
  });

  test.skip(`directly to post: ${post.name}`, async ({ page }) => {
    await page.goto(post.href);
    await page.waitForLoadState("networkidle");

    const actualTitle = await page.title();
    expect(post.name).toEqual(actualTitle);
  });
}
