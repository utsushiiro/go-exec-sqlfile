package main

import "context"

type Post struct {
	ID       int
	AuthorID int
	Content  string
}

func AllPosts(ctx context.Context) ([]Post, error) {
	rows, err := db.QueryContext(ctx, "select * from posts order by id asc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.AuthorID, &post.Content); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
