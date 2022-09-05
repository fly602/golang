package goods

import (
	"fmt"
)

type GoodsInfo struct {
	ID    uint32
	Name  string
	Price float32
	Rest  uint32
}

type Goods struct {
	ID   uint32
	Name string
}

type Shop struct {
	Name  string
	goods map[uint32]*GoodsInfo
}

func InitShop(name string) *Shop {
	return &Shop{
		Name:  name,
		goods: map[uint32]*GoodsInfo{},
	}
}

func (s *Shop) GetGoods(id uint32) *GoodsInfo {
	if s.goods == nil {
		return nil
	}
	fmt.Printf("GetGoods:%+v\n", s.goods[id])
	return s.goods[id]
}

func (s *Shop) PutGoods(g *GoodsInfo) error {
	if s.goods == nil {
		return fmt.Errorf("shop not init")
	}
	for _, good := range s.goods {
		if g.Name == good.Name {
			return fmt.Errorf("goods %s exist", g.Name)
		}
	}
	g.ID = uint32(len(s.goods)) + 1
	s.goods[g.ID] = g
	fmt.Printf("PutGoods:%+v\n", g)
	return nil
}

func (s *Shop) GetGoodsInfo(in *Goods, gi *GoodsInfo) error {
	tmp := s.GetGoods(in.ID)

	if gi == nil {
		return fmt.Errorf("the goods %s not exist", in.Name)
	}
	gi.ID = tmp.ID
	gi.Name = tmp.Name
	gi.Price = tmp.Price
	gi.Rest = tmp.Rest
	fmt.Printf("GetGoodsInfo:%+v\n", gi)
	return nil
}

func (s *Shop) Consume(in *Goods, gi *GoodsInfo) error {
	tmp := s.GetGoods(in.ID)

	if tmp == nil {
		return fmt.Errorf("the goods %s not exist", in.Name)
	}
	if tmp.Rest > 0 {
		tmp.Rest--
		gi.ID = tmp.ID
		gi.Name = tmp.Name
		gi.Price = tmp.Price
		gi.Rest = tmp.Rest
		fmt.Printf("%s sold, id= %d\n", in.Name, in.ID)
		return nil
	} else {
		return fmt.Errorf("%s empty rest", in.Name)
	}
}
