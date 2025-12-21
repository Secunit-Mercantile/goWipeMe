#!/bin/bash
# Helper script to create version tags for releases

set -e

if [ -z "$1" ]; then
    echo "Usage: ./scripts/tag-release.sh <version>"
    echo "Example: ./scripts/tag-release.sh 1.0.0"
    echo ""
    echo "This will create and push tag v1.0.0 which triggers the release workflow"
    exit 1
fi

VERSION=$1

# Remove 'v' prefix if provided
VERSION=${VERSION#v}

# Validate version format (semantic versioning)
if ! [[ $VERSION =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Error: Version must be in format X.Y.Z (e.g., 1.0.0)"
    exit 1
fi

TAG="v$VERSION"

# Check if tag already exists
if git rev-parse "$TAG" >/dev/null 2>&1; then
    echo "Error: Tag $TAG already exists"
    exit 1
fi

echo "Creating release tag: $TAG"
echo ""
echo "Recent commits:"
git log --oneline -5
echo ""

read -p "Create and push tag $TAG? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "Aborted"
    exit 1
fi

# Create annotated tag
git tag -a "$TAG" -m "Release $TAG"

# Push tag to remote
git push origin "$TAG"

echo ""
echo "✓ Tag $TAG created and pushed"
echo "✓ GitHub Actions will now build and create the release"
echo ""
echo "Monitor the workflow at:"
echo "https://github.com/$(git config --get remote.origin.url | sed 's/.*://;s/.git$//')/actions"
