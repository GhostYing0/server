package logic

// 每页显示条数
const PerPage = 10

// Paginator 分页器
type Paginator struct {
	curPage int // 当前页码
	perPage int // 每页条目数

	total     int //总条目数
	totalPage int //总页数
}

func NewPaginator(curPage int, perPage int) *Paginator {
	return &Paginator{curPage: curPage, perPage: perPage}
}

func (this *Paginator) Offset() (offset int) {
	if this.curPage > 1 {
		offset = (this.curPage - 1) * this.perPage
	}
	return
}

func (this *Paginator) SetTotalPage(total int64) {
	this.total = int(total)
	if this.perPage != 0 {
		this.totalPage = this.total / this.perPage
		if this.total%this.perPage != 0 {
			this.totalPage++
		}
	}
}

func (this *Paginator) CurPage() int {
	return this.curPage
}

func (this *Paginator) SetCurPage(curPage int) {
	this.curPage = curPage
}

func (this *Paginator) PerPage() int {
	return this.perPage
}

func (this *Paginator) SetPerPage(perPage int) {
	this.perPage = perPage
}

func (this *Paginator) GetTotal() int {
	return this.total
}

func (this *Paginator) GetTotalPage() int {
	return this.totalPage
}
