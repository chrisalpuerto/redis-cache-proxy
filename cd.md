    deploy:
      needs: ci-checks
      runs-on: ubuntu-latest
      steps:
        - uses: actions/checkout@v3

        # Authenticate GitHub -> GCP
        - name: Authentication to GCP
          uses: google-github-actions/auth@v2
          with:
            credentials_json: ${{ secrets.GCP_CREDENTIALS }}
        - name: Set up Gcloud
          uses: google-github-actions/setup-gcloud@v2
          with:
            project_id: ${{ secrets.GCP_PROJECT_ID }}
          
        # Docker auth for artifact registry
        - name: Configure Docker
          run: gcloud auth configure-docker
        
          # Build and push Docker image for Artifact Registry
        - name: Build and push Docker Image
          run: |
            docker build \
              -t $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE:$GITHUB_SHA .

            docker push \
              $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE:$GITHUB_SHA
              
        # deploy update to cloud run service
        - name: Deploy to Cloud Run
          run: |
            gcloud run deploy $SERVICE_NAME \
              --image $REGION-docker.pkg.dev/$PROJECT_ID/$REPO/$IMAGE:$GITHUB_SHA \
              --region $REGION \
              --platform managed \
              --quiet
