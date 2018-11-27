// Build properties
properties([
  buildDiscarder(
    logRotator(
      artifactDaysToKeepStr: '',
      artifactNumToKeepStr: '',
      daysToKeepStr: '',
      numToKeepStr: '10'
    )
  ),
  disableConcurrentBuilds(),
  disableResume(),
  pipelineTriggers([
    cron('H H * * *')
  ])
])

// Repository name use, must end with / or be '' for none.
// Setting this to '' will also disable any pushing
repository= 'area51/'

// image prefix - project specific
imagePrefix = 'uktransport'

// The image tag (i.e. repository/image but no version)
imageTag=repository + imagePrefix

// The architectures to build, in format recognised by docker
architectures = [ 'amd64', 'arm64v8' ]

node( "AMD64" ) {
  stage('Env') {
    sh 'printenv'
  }
}

// The image version based on the branch name - master branch is latest in docker
version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

// Push an image if we have a repository set
def pushImage = {
  tag -> if( repository != '' ) {
    sh 'docker push ' + tag
  }
}

// Build a service for a specific architecture
def buildArch = {
  nodetag, architecture -> node( nodetag ) {
    stage( "docker" ) {
      checkout scm

      sh './build.sh ' + imageTag + ' ' + architecture + ' ' + version

      pushImage( imageTag + ':' + architecture + '-' + version )
    }
  }
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

      sh './multiarch.sh' +
        ' ' + imageTag +
        ' ' + version +
        ' ' + architectures.join(' ')

      pushImage( imageTag + ':' + version )
    }
  }
}
