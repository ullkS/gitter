import pygame
import random

pygame.init()

# Размеры экрана
screen_width = 800
screen_height = 600

# Цвета
black = (0, 0, 0)
white = (255, 255, 255)
green = (0, 255, 0)
red = (255, 0, 0)
blue = (0, 0, 255)

# Размеры змейки, яблока и камня
snake_size = 20
apple_size = 20
stone_size = 20

# Скорость змейки
snake_speed = 10

game_screen = pygame.display.set_mode((screen_width, screen_height))
pygame.display.set_caption('Змейка')

def show_score(score):
    font = pygame.font.SysFont('Arial', 25)
    score_text = font.render("Счет: " + str(score), True, white)
    game_screen.blit(score_text, [10, 10])

def draw_snake(snake_body):
    for block in snake_body:
        pygame.draw.rect(game_screen, green, [block[0], block[1], snake_size, snake_size])

def generate_apple():
    apple_x = round(random.randrange(0, screen_width - apple_size) / 20) * 20
    apple_y = round(random.randrange(0, screen_height - apple_size) / 20) * 20
    return apple_x, apple_y

def generate_stones(num_stones):
    stones = []
    for _ in range(num_stones):
        stone_x = round(random.randrange(0, screen_width - stone_size) / 20) * 20
        stone_y = round(random.randrange(0, screen_height - stone_size) / 20) * 20
        stones.append((stone_x, stone_y))
    return stones

# Главная функция игры
def game_loop():
    game_over = False
    game_close = False

    # Позиция и направление змейки
    x = screen_width / 2
    y = screen_height / 2
    x_change = 0
    y_change = 0

    snake_body = []
    snake_length = 1

    # Начальное количество камней
    num_stones = 5

    # Генератор первого яблока и камней
    apple_x, apple_y = generate_apple()
    stones = generate_stones(num_stones)

    while not game_over:

        while game_close:
            game_screen.fill(black)
            font = pygame.font.SysFont('Arial', 40)
            game_over_text = font.render("Игра окончена!", True, red)
            restart_text = font.render("Нажмите R для рестарта", True, white)
            game_screen.blit(game_over_text, [screen_width / 2 - game_over_text.get_width() / 2, screen_height / 2 - game_over_text.get_height() / 2])
            game_screen.blit(restart_text, [screen_width / 2 - restart_text.get_width() / 2, screen_height / 2 + restart_text.get_height()])
            pygame.display.update()

            # Обработчик событий при окончании игры
            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    game_over = True
                    game_close = False
                if event.type == pygame.KEYDOWN:
                    if event.key == pygame.K_r:
                        game_loop()

        # Обработчик событий
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                game_over = True
            if event.type == pygame.KEYDOWN:
                if event.key == pygame.K_LEFT:
                    x_change = -snake_size
                    y_change = 0
                elif event.key == pygame.K_RIGHT:
                    x_change = snake_size
                    y_change = 0
                elif event.key == pygame.K_UP:
                    y_change = -snake_size
                    x_change = 0
                elif event.key == pygame.K_DOWN:
                    y_change = snake_size
                    x_change = 0

        # Обновление позиции змейки
        x += x_change
        y += y_change

        # Проверка на выход за границы экрана
        if x >= screen_width:
            x = 0
        elif x < 0:
            x = screen_width - snake_size
        if y >= screen_height:
            y = 0
        elif y < 0:
            y = screen_height - snake_size

        game_screen.fill(black)

        pygame.draw.rect(game_screen, red, [apple_x, apple_y, apple_size, apple_size])

        for stone in stones:
            pygame.draw.rect(game_screen, blue, [stone[0], stone[1], stone_size, stone_size])

        # Обновление змейки
        snake_head = []
        snake_head.append(x)
        snake_head.append(y)
        snake_body.append(snake_head)
        if len(snake_body) > snake_length:
            del snake_body[0]

        # Проверка на столкновение с самой собой
        for block in snake_body[:-1]:
            if block == snake_head:
                game_close = True

        # Проверка на столкновение с камнем
        for stone in stones:
            if x == stone[0] and y == stone[1]:
                game_close = True

        # Отображение змейки
        draw_snake(snake_body)

        # Проверка на столкновение с яблоком
        if x == apple_x and y == apple_y:
            apple_x, apple_y = generate_apple()
            snake_length += 1

        # Отображение счета
        show_score(snake_length - 1)

        pygame.display.update()

        # Ограничение скорости змейки
        pygame.time.Clock().tick(snake_speed)

    pygame.quit()

game_loop()