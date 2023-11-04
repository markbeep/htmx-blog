import { test, expect } from "@playwright/test";

test("index page", async ({ browser }) => {
  const context = await browser.newContext();
  const page = await context.newPage();
  await page.goto("/");
  await page.waitForLoadState("networkidle");
  await expect(page).toHaveTitle(/Hello, I'm Mark/);

  const links = [
    { text: "posts", href: "/posts" },
    { text: "polyring", href: "/polyring" },
    { text: "who is mark?", href: "/about" },
  ];

  for (const link of links) {
    const page = await context.newPage();
    await page.goto("/");

    const linkElement = await page.waitForSelector(
      `a[href="${link.href}"]:has-text("${link.text}")`,
    );
    await linkElement.click();
    await page.waitForLoadState("networkidle");

    await expect(page).toHaveURL(link.href);

    await page.close();
  }
});

test("posts page", async ({ browser }) => {
  const context = await browser.newContext();
  const page = await context.newPage();
  await page.goto("/posts");
  await page.waitForLoadState("networkidle");

  const blogPosts = await page.locator("span").all();
  expect(blogPosts.length > 0).toBeTruthy();
  const postLinks: { name: string; href: string }[] = [];

  const dateRegex =
    /^(Sun|Mon|Tue|Wed|Thu|Fri|Sat)\s(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s\d{2},\s\d{4}/;
  for (const post of blogPosts) {
    const text = await post.textContent();
    expect(text).toMatch(dateRegex);

    const linkElement = post.locator("a");
    const linkText = await linkElement.textContent();

    expect(linkText && linkText.length > 0).toBeTruthy();
    expect(linkElement).toHaveAttribute("href");

    postLinks.push({
      name: linkText ?? "",
      href: (await linkElement.getAttribute("href")) ?? "",
    });
  }

  // Visit every post and check if it works
  for (const post of postLinks) {
    const page = await context.newPage();
    await page.goto("/posts");

    const linkElement = await page.waitForSelector(
      `a[href="${post.href}"]:has-text("${post.name}")`,
    );
    const linkText = await linkElement?.textContent();
    expect(linkText && linkText.length > 0);

    await linkElement.click();
    await page.waitForLoadState("networkidle");

    const actualTitle = await page.title();
    expect(actualTitle).toEqual(linkText);

    await page.close();
  }

  context.close();
});
