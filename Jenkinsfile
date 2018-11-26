// Repository name use, must end with / or be '' for none
repository= 'area51/'
// Disable deployment until refactor is complete
//repository=''

// image prefix
imagePrefix = 'uktransport'

// The image version, master branch is latest in docker
version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

// The architectures to build, in format recognised by docker
architectures = [ 'amd64', 'arm64v8' ]

// Temp docker image name
tempImage = 'temp/' + imagePrefix + ':' + version

// The docker image name
// architecture can be '' for multiarch images
def dockerImage = {
  architecture -> repository + imagePrefix +
    ':' +
    ( architecture=='' ? '' : (architecture + '-') ) +
    version
}

// The multi arch image name
multiImage = repository + imagePrefix + ':' + version

// The go arch
def goarch = {
  architecture -> switch( architecture ) {
    case 'amd64':
      return 'amd64'
    case 'arm32v6':
    case 'arm32v7':
      return 'arm'
    case 'arm64v8':
      return 'arm64'
    default:
      return architecture
  }
}

// goarm is for arm32 only
def goarm = {
  architecture -> switch( architecture ) {
    case 'arm32v6':
      return '6'
    case 'arm32v7':
      return '7'
    default:
      return ''
  }
}

// Build properties
properties([
  buildDiscarder(logRotator(artifactDaysToKeepStr: '', artifactNumToKeepStr: '', daysToKeepStr: '', numToKeepStr: '10')),
  disableConcurrentBuilds(),
  disableResume(),
  pipelineTriggers([
    cron('H H * * *')
  ])
])

// Build a service for a specific architecture
def buildArch = {
  nodetag, architecture ->
    node( nodetag ) {
      stage( "docker" ) {
        checkout scm

        sh 'docker build' +
        ' -t ' + dockerImage( architecture ) +
        ' --build-arg skipTest=true' +
        ' --build-arg arch=' + architecture +
        ' --build-arg goos=linux' +
        ' --build-arg goarch=' + goarch( architecture ) +
        ' --build-arg goarm=' + goarm( architecture ) +
        ' .'

        if( repository != '' ) {
        // Push all built images relevant docker repository
        sh 'docker push ' + dockerImage( architecture )
        } // repository != ''
      }
    }
}

manifests = architectures.collect { architecture -> dockerImage( architecture ) }
manifests = manifests.join(' ')

// Deploy multi-arch image for a service
def multiArchService = {
  tmp ->
    // Create/amend the manifest with our architectures
    sh 'docker manifest create -a ' + multiImage + ' ' + manifests

    // For each architecture annotate them to be correct
    architectures.each {
      architecture -> sh 'docker manifest annotate' +
        ' --os linux' +
        ' --arch ' + goarch( architecture ) +
        ' ' + multiImage +
        ' ' + dockerImage( architecture )
    }

    // Publish the manifest
    sh 'docker manifest push -p ' + multiImage
}

// Build on each platform
parallel (
  'amd64': {
    buildArch( "AMD64", "amd64" )
  },
  'arm64v8': {
    buildArch( "ARM64", "arm64v8" )
  }
)

// The multiarch build only if we have a repository set
if( repository != '' ) {
  node( 'AMD64' ) {
    stage( "Multiarch Image" ) {
      multiArchService( '' )
    }
  }
}
