
# Meal Schedule API
## Введение
Meal Schedule API предоставляет возможность управлять рецептами и расписанием приема пищи.

## Авторизация и аутентификация

- `POST /auth/sign-up`: Регистрация нового пользователя.
    - Тело запроса:
  ```json
  { 
    "name": "<name>",
    "username": "<username>", 
    "password": "<password>" 
  }
  ```
    - Ответ:
  ```json
  { 
    "id": "<user_id>" 
  }
  ```

- `POST /auth/sign-in`: Вход пользователя.
    - Тело запроса:
  ```json
  { 
    "username": "<username>", 
    "password": "<password>" 
  }
  ```
    - Ответ:
  ```json
  { 
    "token": "<token>" 
  }
  ```

## Конечные точки API

### Рецепты

- `POST /api/recipes`: Создание нового рецепта.
    - Тело запроса:
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
    - Ответ:
  ```json
  { 
    "id": "<recipe_id>" 
  }
  ```

- `GET /api/recipes`: Получение всех рецептов.
    - Ответ:
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

- `GET /api/recipes/:id`: Получение рецепта по ID.
    - Ответ:
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

- `PUT /api/recipes/:id`: Обновление рецепта по ID.
    - Тело запроса:
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
    - Ответ:
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

- `DELETE /api/recipes/:id`: Удаление рецепта по ID.
    - Ответ:
  ```json
  { 
    "message": "Recipe deleted successfully." 
  }
  ```

- `GET /api/recipes/public`: Получение общедоступных рецептов.
    - Ответ:
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

- `GET /api/recipes/filter`: Получение отфильтрованных рецептов.
    - Параметры запроса: `?costMoreThan=<cost_more_than>&costLessThan=<cost_less_than>&timeToPrepareMoreThan=<time_to_prepare_more_than>&timeToPrepareLessThan=<time_to_prepare_less_than>&healthyMoreThan=<healthy_more_than>&healthyLessThan=<healthy_less_than>`
    - Ответ:
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

- `GET /api/recipes/userFilter`: Получение отфильтрованных рецептов пользователя.
    - Параметры запроса: `?costMoreThan=<cost_more_than>&costLessThan=<cost_less_than>&timeToPrepareMoreThan=<time_to_prepare_more_than>&timeToPrepareLessThan=<time_to_prepare_less_than>&healthyMoreThan=<healthy_more_than>&healthyLessThan=<healthy_less_than>`
    - Ответ:
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

## Ошибки

В случае ошибки API возвращает ответ в следующем формате:

```json
{
  "error": "<error_message>"
}
```

### Расписание питания

- `POST /api/schedule`: Создание нового элемента расписания

.
- Тело запроса:
  ```json
  { 
    "name": "<meal_name>",
    "at_time": "<at_time>",
    "recipes": ["<recipe_id_1>", "<recipe_id_2>", ...]
  }
  ```
- Ответ:
  ```json
  { 
    "id": "<meal_id>" 
  }
  ```

- `GET /api/schedule/all`: Получение всего расписания.
    - Ответ:
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

- `GET /api/schedule`: Получение расписания за определенный период.
    - Параметры запроса: `?date=<date>&period=<period>`
    - Ответ:
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

- `PUT /api/schedule`: Обновление элемента расписания.
    - Параметры запроса: `?date=<date>`
    - Тело запроса:
  ```json
  { 
    "name": "<meal_name>",
    "at_time": "<at_time>",
    "recipes": ["<recipe_id_1>", "<recipe_id_2>", ...]
  }
  ```
    - Ответ:
  ```json
  { 
    "result": "ok" 
  }
  ```

- `DELETE /api/schedule`: Удаление элемента расписания.
    - Параметры запроса: `?date=<date>`
    - Ответ:
  ```json
  { 
    "result": "ok" 
  }
  ```