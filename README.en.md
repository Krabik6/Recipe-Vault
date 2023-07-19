# Meal Schedule API
# Meal Schedule API

## Table of Contents
1. [Introduction](#introduction)
2. [Authorization and Authentication](#authorization-and-authentication)
3. [API Endpoints](#api-endpoints)
    - [Recipes](#recipes)
    - [Meal Schedule](#meal-schedule)
  4. [Errors](#errors)

## Introduction
The Meal Schedule API provides the ability to manage recipes and meal schedules.

## Authorization and Authentication

- `POST /auth/sign-up`: Register a new user.
    - Request body:
  ```json
  { 
    "name": "<name>",
    "username": "<username>", 
    "password": "<password>" 
  }
  ```
    - Response:
  ```json
  { 
    "id": "<user_id>" 
  }
  ```

- `POST /auth/sign-in`: User login.
    - Request body:
  ```json
  { 
    "username": "<username>", 
    "password": "<password>" 
  }
  ```
    - Response:
  ```json
  { 
    "token": "<token>" 
  }
  ```

## API Endpoints

### Recipes

- `POST /api/recipes`: Create a new recipe.
    - Request body:
  ```json
  { 
    "title": "<recipe_title>", 
    "description": "<recipe_description>", 
    "public": "<is_public>",
    "cost": "<cost>",
    "timeToPrepare": "<time_to_prepare>",
    "healthy": "<healthy>",
    "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
    "ingredients": [
      {
        "name": "<ingredient_name>",
        "quantity": "<quantity>",
        "unit": "<unit>"
      },
      ...
    ]
  }
  ```
    - Response:
  ```json
  { 
    "id": "<recipe_id>" 
  }
  ```

- `GET /api/recipes`: Get all recipes.
    - Response:
  ```json
  [ 
    { 
      "id": "<recipe_id>", 
      "title": "<recipe_title>", 
      "description": "<recipe_description>", 
      "public": "<is_public>",
      "cost": "<cost>",
      "timeToPrepare": "<time_to_prepare>",
      "healthy": "<healthy>",
      "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
      "ingredients": [
        {
          "id": "<ingredient_id>",
          "name": "<ingredient_name>",
          "price": "<price>",
          "unit": "<unit>",
          "unitShort": "<unit_short>",
          "unitLong": "<unit_long>",
          "possibleUnits": ["<possible_unit_1>", "<possible_unit_2>", ...],
          "protein": "<protein>",
          "fat": "<fat>",
          "carbs": "<carbs>",
          "aisle": "<aisle>",
          "image": "<image>",
          "categoryPath": ["<category_path_1>", "<category_path_2>", ...],
          "consistency": "<consistency>",
          "external_id": "<external_id>",
          "amount": "<amount>"
        },
        ...
      ]
    }, 
    ... 
  ]
  ```

- `GET /api/recipes/:id`: Get a recipe by ID.
    - Response:
  ```json
  { 
    "id": "<recipe_id>", 
    "title": "<recipe_title>", 
    "description": "<recipe_description>", 
    "public": "<is_public>",
    "cost": "<cost>",
    "timeToPrepare": "<time_to_prepare>",
    "healthy": "<healthy>",
    "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
    "ingredients": [
      {
        "id": "<ingredient_id>",
        "name": "<ingredient_name>",
        "price": "<price>",
        "unit": "<unit>",
        "unitShort": "<unit_short>",
        "unitLong": "<unit_long>",
        "possibleUnits": ["<possible_unit_1>", "<possible_unit_2>", ...],
        "protein": "<protein>",
        "fat": "<fat>",
        "carbs": "<carbs>",
        "aisle": "<aisle>",
        "image": "<image>",
        "categoryPath": ["<category_path_1>", "<category_path_2>", ...],
        "consistency": "<consistency>",
        "external_id": "<external_id>",
        "amount": "<amount>"
      },
      ...
    ]
  }
  ```

- `PUT /api/recipes/:id`: Update a recipe by ID.
    - Request body:
  ```json
  { 
    "title": "<recipe_title>", 
    "description": "<recipe_description>", 
    "public": "<is_public>",
    "cost": "<cost>",
    "timeToPrepare": "<time_to_prepare>",
    "healthy": "<healthy>",
    "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
    "ingredients": [
      {
        "name": "<ingredient_name>",
        "quantity": "<quantity>",
        "unit": "<unit>"
      },
      ...
    ]
  }
  ```
    - Response:
  ```json
  { 
    "id": "<recipe_id>", 
    "title": "<recipe_title>", 
    "description": "<recipe_description>", 
    "public": "<is_public>",
    "cost": "<cost>",
    "timeToPrepare": "<time_to_prepare>",
    "healthy": "<healthy>",
    "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
    "ingredients": [
      {
        "id": "<ingredient_id>",
        "name": "<ingredient_name>",
        "price": "<price>",
        "unit": "<unit>",
        "unitShort": "<unit_short>",
        "unitLong": "<unit_long>",
        "possibleUnits": ["<possible_unit_1>", "<possible_unit_2>", ...],
        "protein": "<protein>",
        "fat": "<fat>",
        "carbs": "<carbs>",
        "aisle": "<aisle>",
        "image": "<image>",
        "categoryPath": ["<category_path_1>", "<category_path_2>", ...],
        "consistency": "<consistency>",
        "external_id": "<external_id>",
        "amount": "<amount>"
      },
      ...
    ]
  }
  ```

- `DELETE /api/recipes/:id`: Delete a recipe by ID.
    - Response:
  ```json
  { 
    "message": "Recipe deleted successfully." 
  }
  ```

- `GET /api/recipes/public`: Get public recipes.
    - Response:
  ```json
  [ 
    { 
      "id": "<recipe_id>", 
      "title": "<recipe_title>", 
      "description": "<recipe_description>", 
      "public": "<is_public>",
      "cost": "<cost>",
      "timeToPrepare": "<time_to_prepare>",
      "healthy": "<healthy>",
      "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
      "ingredients": [
        {
          "id": "<ingredient_id>",
          "name": "<ingredient_name>",
          "price": "<price>",
          "unit": "<unit>",
          "unitShort": "<unit_short>",
          "unitLong": "<un  it_long>",
          "possibleUnits": ["<possible_unit_1>", "<possible_unit_2>", ...],
          "protein": "<protein>", 
          "fat": "<fat>",
          "carbs": "<carbs>",
          "aisle": "<aisle>",
          "image": "<image>",
          "categoryPath": ["<category_path_1>", "<category_path_2>", ...],
          "consistency": "<consistency>",
          "external_id": "<external_id>",
          "amount": "<amount>"
          },
          ...
          ]
]
}
  ```

- `GET /api/recipes/filter`: Get filtered recipes.
    - Request parameters: `?costMoreThan=<cost_more_than>&costLessThan=<cost_less_than>&timeToPrepareMoreThan=<time_to_prepare_more_than>&timeToPrepareLessThan=<time_to_prepare_less_than>&healthyMoreThan=<healthy_more_than>&healthyLessThan=<healthy_less_than>`
    - Response:
  ```json
  [ 
    { 
      "id": "<recipe_id>", 
      "title": "<recipe_title>", 
      "description": "<recipe_description>", 
      "public": "<is_public>",
      "cost": "<cost>",
      "timeToPrepare": "<time_to_prepare>",
      "healthy": "<healthy>",
      "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
      "ingredients": [
        {
          "id": "<ingredient_id>",
          "name": "<ingredient_name>",
          "price": "<price>",
          "unit": "<unit>",
          "unitShort": "<unit_short>",
          "unitLong": "<unit_long>",
          "possibleUnits": ["<possible_unit_1>", "<possible_unit_2>", ...],
          "protein": "<protein>",
          "fat": "<fat>",
          "carbs": "<carbs>",
          "aisle": "<aisle>",
          "image": "<image>",
          "categoryPath": ["<category_path_1>", "<category_path_2>", ...],
          "consistency": "<consistency>",
          "external_id": "<external_id>",
          "amount": "<amount>"
        },
        ...
      ]
    }, 
    ... 
  ]
  ```

- `GET /api/recipes/userFilter`: Get user's filtered recipes.
    - Request parameters: `?costMoreThan=<cost_more_than>&costLessThan=<cost_less_than>&timeToPrepareMoreThan=<time_to_prepare_more_than>&timeToPrepareLessThan=<time_to_prepare_less_than>&healthyMoreThan=<healthy_more_than>&healthyLessThan=<healthy_less_than>`
    - Response:
  ```json
  [ 
    { 
      "id": "<recipe_id>", 
      "title": "<recipe_title>", 
      "description": "<recipe_description>", 
      "public": "<is_public>",
      "cost": "<cost>",
      "timeToPrepare": "<time_to_prepare>",
      "healthy": "<healthy>",
      "imageURLs": ["<image_url_1>", "<image_url_2>", ...],
      "ingredients": [
        {
          "id": "<ingredient_id>",
          "name": "<ingredient_name>",
          "price": "<price>",
          "unit": "<unit>",
          "unitShort": "<unit_short>",
          "unitLong": "<unit_long>",
          "possibleUnits": ["<possible_unit_1>", "<possible_unit_2>", ...],
          "protein": "<protein>",
          "fat": "<fat>",
          "carbs": "<carbs>",
          "aisle": "<aisle>",
          "image": "<image>",
          "categoryPath": ["<category_path_1>", "<category_path_2>", ...],
          "consistency": "<consistency>",
          "external_id": "<external_id>",
          "amount": "<amount>"
        },
        ...
      ]
    }, 
    ... 
  ]
  ```

## Errors

In case of an error, the API returns a response in the following format:

```json
{
  "error": "<error_message>"
}
```

### Meal Schedule

- `POST /api/schedule`: Create a new schedule item.
    - Request body:
  ```json
  { 
    "name": "<meal_name>",
    "at_time": "<at_time>",
    "recipes": ["<recipe_id_1>", "<recipe_id_2>", ...]
  }
  ```
    - Response:
  ```json
  { 
    "id": "<meal_id>" 
  }
  ```

- `GET /api/schedule/all`: Get the entire schedule.
    - Response:
  ```json
  [ 
    { 
      "id": "<meal_id>", 
      "name": "<meal_name>",
      "at_time": "<at_time>",
      "recipes": ["<recipe_id_1>", "<recipe_id_2>", ...]
    }, 
    ... 
  ]
  ```

- `GET /api/schedule`: Get the schedule for a specific period.
    - Request parameters: `?date=<date>&period=<period>`
    - Response:
  ```json
  [ 
    { 
      "id": "<meal_id>", 
      "name": "<meal_name>",
      "at_time": "<at_time>",
      "recipes": ["<recipe_id_1>", "<recipe_id_2>", ...],
      "date": "<date>", 
      "period": "<period>" 
    }, 
    ... 
  ]
  ```

- `PUT /api/schedule`: Update a schedule item.
    - Request parameters: `?date=<date>`
    - Request body:
  ```json
  { 
    "name": "<meal_name>",
    "at_time": "<at_time>",
    "recipes": ["<recipe_id_1>", "<recipe_id_2>", ...]
  }
  ```
    - Response:
  ```json
  { 
    "result": "ok" 
  }
  ```

- `DELETE /api/schedule`: Delete a schedule item.
    - Request parameters: `?date=<date>`
    - Response:
  ```json
  { 
    "result": "ok" 
  }
  ```