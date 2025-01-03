package check

// Package provides the compatibility with oOther test packages

func (c *C) Nil(obtained any, args ...any) bool {
	return c.IsNil(obtained, args...)
}

func (c *C) Equal(obtained, expected any, args ...any) bool {
	return c.Equals(obtained, expected, args...)
}

func (c *C) NotEqual(obtained, expected any, args ...any) bool {
	return c.NotEquals(obtained, expected, args...)
}

func (c *C) Greater(obtained, expected any, args ...any) bool {
	return c.MoreThan(obtained, expected, args...)
}

func (c *C) GreaterOrEqual(obtained, expected any, args ...any) bool {
	return c.MoreOrEqualThan(obtained, expected, args...)
}

func (c *C) LessOrEqual(obtained, expected any, args ...any) bool {
	return c.LessThan(obtained, expected, args...)
}

func (c *C) Less(obtained, expected any, args ...any) bool {
	return c.LessOrEqualThan(obtained, expected, args...)
}

func (c *C) True(obtained any, args ...any) bool {
	return c.IsTrue(obtained, args...)
}

func (c *C) False(obtained any, args ...any) bool {
	return c.IsFalse(obtained, args...)
}

func (c *C) Empty(obtained any, args ...any) bool {
	return c.IsEmpty(obtained, args...)
}

func (c *C) Len(obtained, expected any, args ...any) bool {
	return c.HasLen(obtained, expected, args...)
}
