import urllib.request, urllib.parse, json, pandas as pd, datetime

# ---- CONFIGURATION ----
GITHUB_TOKEN = "xxxx"   # use personal access token
ORG = "xxxx"               # org or username
FROM_DATE = "2024-10-01T00:00:00Z"
TO_DATE = "2025-12-01T00:00:00Z"
HEADERS = {"Authorization": f"token {GITHUB_TOKEN}",
           "Accept": "application/vnd.github+json"}

# ---- BASIC GET REQUEST HANDLER WITH PAGINATION ----
def github_get(url):
    results = []
    while url:
        req = urllib.request.Request(url, headers=HEADERS)
        with urllib.request.urlopen(req) as resp:
            data = json.load(resp)
            results.extend(data if isinstance(data, list) else [data])
            # handle pagination via 'Link' header if present
            links = resp.headers.get("Link", "")
            next_url = None
            for link in links.split(","):
                if 'rel="next"' in link:
                    next_url = link[link.find("<")+1:link.find(">")]
            url = next_url
    return results

# ---- FETCH REPOSITORIES ----
repos = github_get(f"https://api.github.com/users/{ORG}/repos?per_page=2")

print("repos: count =", len(repos))
target_repos = {}
# print(repos)
for r in repos:
    if r["name"] == "experiments":
        target_repos[r["name"]] = r

print("Using repos:")
for r in target_repos.values():
    print(" -", r["name"])

commit_data, pr_comment_data = [], []

# ---- FETCH COMMITS & STATS ----
for repo in target_repos.values():
    name = repo["name"]
    commits_url = f"https://api.github.com/repos/{ORG}/{name}/commits?since={FROM_DATE}&until={TO_DATE}&per_page=100"
    commits = github_get(commits_url)
    print(f"Repo '{name}': commits count =", len(commits))

    for c in commits:
        if not c.get("author"):  # skip unknown author
            continue
        user = c["author"]["login"]
        sha = c["sha"]
        # get commit stats (additions/deletions)
        commit_detail_url = f"https://api.github.com/repos/{ORG}/{name}/commits/{sha}"
        detail = github_get(commit_detail_url)[0]
        stats = detail.get("stats", {})
        commit_data.append({
            "repo": name,
            "user": user,
            "sha": sha,
            "additions": stats.get("additions", 0),
            "deletions": stats.get("deletions", 0),
            "total": stats.get("total", 0),
        })

    # ---- PULL REQUEST COMMENTS ----
    comments_url = f"https://api.github.com/repos/{ORG}/{name}/pulls/comments?since={FROM_DATE}&until={TO_DATE}&per_page=100"
    print(f"Fetching PR comments for repo '{name}' from {comments_url}")
    comments = github_get(comments_url)
    for c in comments:
        if not c.get("user"):
            continue
        pr_comment_data.append({
            "repo": name,
            "user": c["user"]["login"],
        })

print("Total commits fetched:", len(commit_data))
# print("Commit data sample:", commit_data)

print("Total PR comments fetched:", len(pr_comment_data))
print("PR comment data sample:", pr_comment_data)

# ---- AGGREGATE USING PANDAS ----
df_commits = pd.DataFrame(commit_data)
df_comments = pd.DataFrame(pr_comment_data)

print("df_commits: count =", len(df_commits))
# print(df_commits.to_string(index=False))
print("\ndf_comments: count =", len(df_comments))
# print(df_comments.to_string(index=False))

exit()

# lines of code metrics
repo_user_stats = df_commits.groupby(["repo", "user"]).sum().reset_index()
overall_user_stats = df_commits.groupby("user").sum().reset_index()

# commits per user/repo
commit_counts = df_commits.groupby(["repo", "user"]).size().reset_index(name="commit_count")
overall_commits = df_commits.groupby("user").size().reset_index(name="commit_count")

# PR comments per user
pr_comments = df_comments.groupby(["repo", "user"]).size().reset_index(name="pr_comments")
overall_comments = df_comments.groupby("user").size().reset_index(name="pr_comments")

# ---- EXPORT TO CSV ----
repo_user_stats.to_csv("repo_user_loc_metrics.csv", index=False)
overall_user_stats.to_csv("overall_user_loc_metrics.csv", index=False)
commit_counts.to_csv("repo_user_commit_counts.csv", index=False)
overall_commits.to_csv("overall_user_commit_counts.csv", index=False)
pr_comments.to_csv("repo_user_pr_comments.csv", index=False)
overall_comments.to_csv("overall_user_pr_comments.csv", index=False)

print("✅ CSV reports generated successfully.")
