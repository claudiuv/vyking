# vyking

This application is a starting point for an iGaming platform.
The following steps are needed for the app database to be deployed in a docker containr:
    a. edit file docker-compose.yml for passwords
    b. run docker-compose up -d
    c. run docker exec -i mysql-container mysql -u root -p strong_password --database=Vyking < ./script.sql
        with strong_password replaced with the password set in docker-compose.yml for root user

This solution contains the following api endpoints:
- /api/v1/distributePrize [post] - distribute the prizes for a tournament based on playes rankings
- /api/v1/rankings [get] - get the players ranking list

