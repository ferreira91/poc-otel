# Use an official openjdk image as the base image
FROM openjdk:17-jdk-slim

# Set the working directory within the container
WORKDIR /otel

# Copy the compiled jar file from the build environment to the container
COPY build/libs/kotlin-0.0.1-SNAPSHOT.jar /otel/kotlin-otel.jar

# Specify the command to run when the container starts
CMD ["java", "-jar", "kotlin-otel.jar"]