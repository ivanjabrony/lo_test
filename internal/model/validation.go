package model

import "errors"

func ValidateTask(task Task) error {
	if task.Id < 0 {
		return errors.New("invalid taskId in task: negative values are forbidden")
	}
	if task.Status == "" {
		return errors.New("invalid status in task: empty status is forbidden")
	}
	if task.Status != Done && task.Status != InProgress && task.Status != Created {
		return errors.New("invalid status in task: unknown type")
	}

	return nil
}

func ValidateFilter(filter Filter) error {
	if filter.Status != Done && filter.Status != InProgress && filter.Status != Created {
		return errors.New("invalid status in filter: unknown type")
	}
	return nil
}
