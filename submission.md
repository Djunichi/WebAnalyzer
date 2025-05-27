
# WebAnalyzer Submission Document

## Overview

This Go-based web application analyzes a given webpage by URL and provides the following insights:
- HTML version detection
- Page title
- Number of heading tags (h1–h6)
- Count of internal, external, and inaccessible links
- Presence of a login form (`<input type="password">`)

---

## Assumptions & Implementation Notes

- **HTML version** is derived from the `<!DOCTYPE>` declaration.
- **Login form detection** is implemented by searching for `<input type="password">` elements.
- Internal links are those sharing the same hostname as the input URL; others are external.
- Inaccessible links are identified using `http.Head()` requests.
- All parsed analyses are stored in a PostgreSQL database and can be reopened by ID.
- Basic frontend is implemented in static HTML with JavaScript form submission.

---

## Suggestions for Improvement

1. **Test Coverage**:
    - Extend unit tests to cover analysis and handler logic (currently only helper functions are tested).

2. **Frontend Enhancements**:
    - Improve UI for a more polished experience.

3. **History View**:
    - Add a dedicated UI page to list all past analyses instead of accessing them by ID.

4. **CI/CD**:
    - Add automated testing and linting via GitHub Actions or GitLab CI pipelines and deployment to cloud infrastructure.

5. **Caching**:
   -  Add analysis results caching to decrease server requests amount. Move in-memory cache to external (redis etc.) to support multiple application instances.
   
6. **Monitoring**
   -  Add external logging and monitoring tools.

