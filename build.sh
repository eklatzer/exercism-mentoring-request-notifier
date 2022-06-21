#! /bin/bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

platforms=(
    "darwin/386"
    "darwin/amd64"
    "linux/386"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "windows/386"
    "windows/amd64"
    "windows/arm"
)

echo "> re-creating bin-folder"
rm -r bin
mkdir bin

echo "> building package ${package}"

for platform in "${platforms[@]}"; do
  echo " > building for ${platform}"
  IFS='/' read -r -a platform_split <<<"${platform}"
  GOOS=${platform_split[0]}
  GOARCH=${platform_split[1]}
  output_name=$package'_'$GOOS'_'$GOARCH

  if [ "${GOOS}" = "windows" ]; then
     output_name+='.exe'
  fi


  if ! env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "bin/${output_name}" "${package}"; then
    echo "  > [ERROR]: failed to build for ${platform}"
    exit 1
  fi
done

echo "> finished build"
