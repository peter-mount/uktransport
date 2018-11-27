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

// The image version based on the branch name - master branch is latest in docker
version=BRANCH_NAME
if( version == 'master' ) {
  version = 'latest'
}

// ======================================================================
// Do not modify anything below this point
// ======================================================================

// Unused for jobs that don't generate artifacts but get the
// JOB_NAME consists of the job name but if it ends with the branch name then
// remove that to get the full name
jobname=JOB_NAME.split('/')
if jobname[-1] == BRANCH_NAME {
  jobname.removeLast()
}

// The artifact path in our repository
repoPath='https://nexus.area51.onl/repository/snapshots/' + jobname.join('/')
// The file name sans any file suffix, contains job name, branch name & build number
repoName=jobname[-1] + '-' + BRANCH_NAME + '.' + BUILD_NUMBER

// Push an image if we have a repository set
def pushImage = {
  tag -> if( repository != '' ) {
    sh 'docker push ' + tag
  }
}

// Build a service for a specific architecture
def buildArch = {
  nodetag, architecture -> node( nodetag ) {
    withEnv([
      'UPLOAD_PATH=' + repoPath,
      'UPLOAD_NAME=' + repoName
    ]) {
      stage( "docker" ) {
        checkout scm

        sh './build.sh ' + imageTag + ' ' + architecture + ' ' + version

        pushImage( imageTag + ':' + architecture + '-' + version )
      }
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
