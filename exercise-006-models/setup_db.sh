#!/bin/bash

# Set the database name
DB_NAME="exercise_006"

# Create the database
psql -U elijah -d postgres -c "CREATE DATABASE $DB_NAME;"

# Print success message
echo "Database '$DB_NAME' created successfully."
