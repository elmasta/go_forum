package forum

import models "forum/pkg/models"

// HasPost checks if a given post (b) exists in a slice of posts (a).
// It iterates over the slice of posts and compares the ID of each post with the ID of the given post.
// If a match is found, it returns true. If no match is found after checking all posts, it returns false.
func HasPost(a []models.Post, b models.Post) bool {
	for _, p := range a {
		if p.ID == b.ID {
			return true
		}
	}
	return false
}
