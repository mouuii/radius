import radius as radius

@description('Specifies the location for resources.')
param location string = 'global'

@description('Specifies the image of the container resource.')
param image string = 'radiusdev.azurecr.io/magpiego:latest'

@description('Specifies the environment for resources.')
param environment string

resource app 'Applications.Core/applications@2022-03-15-privatepreview' = {
  name: 'corerp-resources-container-httproute'
  location: location
  properties: {
    environment: environment
  }
}

resource backend 'Applications.Core/containers@2022-03-15-privatepreview' = {
  name: 'backend'
  location: location
  properties: {
    application: app.id
    container: {
      image: 'jkotalik.azurecr.io/backend:latest'
      ports: {
        web: {
          containerPort: 80
          provides: backendhttp.id
        }
      }
      volumes:{
        'my-volume':{
          kind: 'ephemeral'
          mountPath:'/tmpfs'
          managedStore:'memory'
        }
      }
    }
    connections: {}
  }
}

resource backendhttp 'Applications.Core/httpRoutes@2022-03-15-privatepreview' = {
  name: 'backend'
  location: location
  properties: {
    application: app.id
  }
}

resource frontend 'Applications.Core/containers@2022-03-15-privatepreview' = {
  name: 'frontend'
  location: location
  properties: {
    application: app.id
    container: {
      image: 'jkotalik.azurecr.io/frontend:latest'
      ports: {
        web: {
          containerPort: 80
          provides: frontendhttp.id
        }
      }
      env: {
        // for this example, populate values with existing stuff still populating 
        CONNECTION__BACKEND__HOSTNAME: backendhttp.properties['hostname']
        CONNECTION__BACKEND__PORT: '${backendhttp.properties.port}'
      }
    }
    connections: {
      backend: {
        source: backendhttp.id
      }
    }
  }
}

resource frontendhttp 'Applications.Core/httpRoutes@2022-03-15-privatepreview' = {
  name: 'frontend'
  location: location
  properties: {
    application: app.id
  }
}
