docker run -d --name matchflow -p8080:80 -v /root/matchflow/apps:/services/matchflow/apps -v /root/matchflow/conf:/services/matchflow/conf wangzhsh/matchflow:0.1.0 