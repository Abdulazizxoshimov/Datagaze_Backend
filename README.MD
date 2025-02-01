# Weather Fetch and Store Application

This project is a weather data fetching and storing application that retrieves weather data for specified cities and stores it in a database. The application uses Go (Golang) as the main programming language and supports handling errors and logging efficiently.

## Features

- Fetches weather data for a given city using a weather API.
- Stores the fetched data into a database.
- Handles errors gracefully and logs relevant information.
- Provides an HTTP API for fetching weather information.

## Technologies Used

- **Go (Golang)**: The main programming language.
- **Gin**: A fast web framework for Go.
- **PostgreSQL**: A relational database for storing weather data.
- **Logger**: Custom logging for error and info tracking.

## Setup

### Prerequisites

Before running the project, make sure you have the following installed:

- Go (Golang) [Download Go](https://golang.org/dl/)
- PostgreSQL [Download PostgreSQL](https://www.postgresql.org/download/)
- Git [Download Git](https://git-scm.com/)

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/weather-fetch-store.git
   cd weather-fetch-store
