Pull Request Review Comments: These are comments made on specific lines of code or files within a pull request review.
Endpoint: GET /repos/{owner}/{repo}/pulls/{pull_number}/comments

Issue Comments (Conversation Comments): These are general comments made in the "Conversation" tab of a pull request, which GitHub treats as a type of issue.
Endpoint: GET /repos/{owner}/{repo}/issues/{issue_number}/comments

Pull Request Reviews: These are entire reviews, which can include multiple pull request review comments and an overall review body.
Endpoint: GET /repos/{owner}/{repo}/pulls/{pull_number}/reviews
