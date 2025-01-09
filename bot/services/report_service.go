package services

import (
	"fmt"
	"gorag-telegram-bot/database"
	"gorag-telegram-bot/models"
	"time"
)

// Количество вопросов за разные периоды
type LikedStats struct {
	Liked *bool
	Count int64
}

type DateValues struct {
	Date  time.Time
	Count int
}

// Получение пользователя по TelegramID
func GetCommonReport() string {
	// Количество подписчиков
	var totalUsers int64
	database.DB.Model(&models.User{}).Count(&totalUsers)

	// Получаем текущее время и начало/конец необходимых периодов
	now := time.Now()
	startOfToday := now.Truncate(24 * time.Hour)
	startOfYesterday := startOfToday.AddDate(0, 0, -1)
	startOfThisMonth := now.AddDate(0, 0, -now.Day()+1)
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)
	endOfYesterday := startOfToday
	endOfLastMonth := startOfThisMonth

	// Количество подписчиков за разные периоды
	var newUsersPastMonth, newUsersThisMonth, newUsersYesterday, newUsersToday int64
	database.DB.Model(&models.User{}).Where("connection_date >= ? AND connection_date < ?", startOfLastMonth, endOfLastMonth).Count(&newUsersPastMonth)
	database.DB.Model(&models.User{}).Where("connection_date >= ?", startOfThisMonth).Count(&newUsersThisMonth)
	database.DB.Model(&models.User{}).Where("connection_date >= ? AND connection_date < ?", startOfYesterday, endOfYesterday).Count(&newUsersYesterday)
	database.DB.Model(&models.User{}).Where("connection_date >= ?", startOfToday).Count(&newUsersToday)

	// Получаем статистику по группам
	var (
		statsPastMonth, statsThisMonth, statsYesterday, statsToday []LikedStats
	)

	// Подсчеты для прошлого месяца
	database.DB.Model(&models.Message{}).
		Select("liked, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startOfLastMonth, endOfLastMonth).
		Group("liked").
		Scan(&statsPastMonth)

	// Подсчеты для текущего месяца
	database.DB.Model(&models.Message{}).
		Select("liked, COUNT(*) as count").
		Where("created_at >= ?", startOfThisMonth).
		Group("liked").
		Scan(&statsThisMonth)

	// Подсчеты для вчера
	database.DB.Model(&models.Message{}).
		Select("liked, COUNT(*) as count").
		Where("created_at >= ? AND created_at < ?", startOfYesterday, endOfYesterday).
		Group("liked").
		Scan(&statsYesterday)

	// Подсчеты для сегодня
	database.DB.Model(&models.Message{}).
		Select("liked, COUNT(*) as count").
		Where("created_at >= ?", startOfToday).
		Group("liked").
		Scan(&statsToday)

	// Подсчет итоговых значений
	totalQuestions, totalLikes, totalDislikes := extractLikeDislikeCounts(append(append(append(statsPastMonth, statsThisMonth...), statsYesterday...), statsToday...))

	questionsPastMonth, likesPastMonth, dislikesPastMonth := extractLikeDislikeCounts(statsPastMonth)
	questionsThisMonth, likesThisMonth, dislikesThisMonth := extractLikeDislikeCounts(statsThisMonth)
	questionsYesterday, likesYesterday, dislikesYesterday := extractLikeDislikeCounts(statsYesterday)
	questionsToday, likesToday, dislikesToday := extractLikeDislikeCounts(statsToday)

	// Формируем отчет
	report := fmt.Sprintf(
		"Подписчиков всего: %d\n"+
			"Новые подписчики:\n"+
			" - Прошлый месяц: %d\n"+
			" - Текущий месяц: %d\n"+
			" - Вчера: %d\n"+
			" - Сегодня: %d\n\n"+
			"Вопросов всего: %d (👍%d / 👎%d)\n"+
			"Вопросы:\n"+
			" - Прошлый месяц: %d (👍%d / 👎%d)\n"+
			" - Текущий месяц: %d (👍%d / 👎%d)\n"+
			" - Вчера: %d (👍%d / 👎%d)\n"+
			" - Сегодня: %d (👍%d / 👎%d)\n",
		totalUsers, newUsersPastMonth, newUsersThisMonth, newUsersYesterday, newUsersToday,
		totalQuestions, totalLikes, totalDislikes,
		questionsPastMonth, likesPastMonth, dislikesPastMonth,
		questionsThisMonth, likesThisMonth, dislikesThisMonth,
		questionsYesterday, likesYesterday, dislikesYesterday,
		questionsToday, likesToday, dislikesToday,
	)
	return report
}

func extractLikeDislikeCounts(stats []LikedStats) (total, likes, dislikes int64) {
	for _, stat := range stats {
		total += stat.Count
		if stat.Liked != nil {
			if *stat.Liked {
				likes += stat.Count
			} else {
				dislikes += stat.Count
			}
		}
	}
	return
}

func GetUserRegistrationsTimeline(startFrom time.Time) []DateValues {
	// Получение данных
	var subscriptionData []DateValues

	database.DB.Raw(`
WITH dates AS (SELECT generate_series(
                              (SELECT MIN(created_at) FROM users),
                              NOW() + INTERVAL '1 day',
                              INTERVAL '1 day'
                      )::date AS date)
SELECT dates.date,
       COUNT(users.id) AS count
FROM dates
         LEFT JOIN
     users
     ON
         users.created_at <= dates.date
             AND (users.is_admin = false)
GROUP BY dates.date
HAVING dates.date > ?
ORDER BY dates.date
    `, startFrom).Scan(&subscriptionData)
	return subscriptionData
}

func GetSubscriptionsTimeline(startFrom time.Time) []DateValues {
	// Получение данных
	var subscriptionData []DateValues

	database.DB.Raw(`
        WITH dates AS (
    SELECT generate_series(
                   (SELECT MIN(created_at) FROM subscriptions),
                   NOW() + INTERVAL '1 day',
                   INTERVAL '1 day'
           )::date AS date
)
SELECT
    dates.date,
    COUNT(subscriptions.id) AS count
FROM
    dates
        LEFT JOIN
    subscriptions
    ON
        subscriptions.created_at <= dates.date
            AND (subscriptions.expires_at > dates.date OR subscriptions.expires_at IS NULL)
GROUP BY
    dates.date
HAVING dates.date > ?
ORDER BY
    dates.date
    `, startFrom).Scan(&subscriptionData)
	return subscriptionData
}
