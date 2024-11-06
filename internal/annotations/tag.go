package annotations

// * Tag

type (
	Tag struct {
		aChan aChan
		value string
	}
)

func (a *Annotation) Tag() (at *Tag) {
	at = &Tag{
		aChan: make(aChan),
	}

	if v, ok := a.annotations[string(KeyTag)]; ok {
		at.value = v
	}

	go func() {
		for {
			select {
			case x := <-at.aChan:
				if a.annotations == nil {
					a.annotations = make(map[string]string)
				}

				a.annotations[string(x.key)] = x.value
			case <-a.ctx.Done():
				return
			}
		}
	}()

	return at
}

func (a *Tag) Get() string {
	return a.value
}

func (a *Tag) Set(tag string) {
	a.aChan.Send(KeyTag, tag)
}

func (a *Tag) IsNull() bool {
	return a.value == ""
}
