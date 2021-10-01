set -e

version=$(cut -d'=' -f2- .release)
if [[ -z ${version} ]]; then
    echo "Invalid version set in .release"
    exit 1
fi


if [[ -z ${GITHUB_TOKEN} ]]; then
    echo "GITHUB_TOKEN not set. Usage: GITHUB_TOKEN=<TOKEN> ./hack/release.sh"
    exit 1
fi

echo "Publishing release ${version}"

generate_changelog() {
    local version=$1

    # generate changelog from github
    git-chglog --output CHANGELOG.md
    sed -i ''  '$d' CHANGELOG.md
}

update_chart_yamls() {
    local version=$1

    sed -i ''  "s/version.*/version: ${version}/" helm/kubediff/Chart.yaml
    sed -i ''  "s/appVersion.*/appVersion: ${version}/" helm/kubediff/Chart.yaml
    sed -i ''  "s/\btag:.*/tag: ${version}/" helm/kubediff/values.yaml
    sed -i ''  "s/\bimage: \"arriqaaq\/kubediff.*\b/image: \"arriqaaq\/kubediff:${version}/g" hack/deploy.yaml
}

publish_release() {
    local version=$1

    github-release release \
	   --user arriqaaq \
	   --repo kubediff \
	   --tag $version \
	   --name "$version" \
	   --description "$version"
}

update_chart_yamls $version
# generate_changelog $version
make release
publish_release $version

echo "Release ${version} published."
