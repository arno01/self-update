# Hints

Setup and configure `gh` tool.

```
sudo apt -y install gh
```

```
gh auth login
```

## Remove all releases

```
gh release list | awk '{print $1}' | xargs -n1 gh release delete -y
```

## Remove all tags

```
git fetch
git tag -l | xargs -n 1 git push --delete origin
git tag | xargs git tag -d
```

## Remove all history

Except for the PR's / branches ..

```
rm -rf .git
git init
git branch -m main
git add .
git commit -S -m 'first commit'
git tag v0.0.1
git remote add origin git@github.com:arno01/self-update.git
git push -f origin main
```

