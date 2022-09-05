package goods

import (
	"context"
	"fmt"
	"go-community/grpc/sell/proto/sell"

	"github.com/golang/protobuf/ptypes/empty"
)

type GoodsInfo struct {
	ID    uint32
	Name  string
	Price float32
	Rest  uint32
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
	return nil
}

func (s *Shop) GetGoodsInfo(ctx context.Context, in *sell.Goods) (*sell.GoodsInfo, error) {
	gi := s.GetGoods(in.Id)

	if gi == nil {
		return nil, fmt.Errorf("the goods %s not exist", in.Name)
	}
	goods := &sell.GoodsInfo{
		G: &sell.Goods{
			Id:   gi.ID,
			Name: gi.Name,
		},
		Price: gi.Price,
		Rest:  gi.Rest,
	}
	return goods, nil
}

func (s *Shop) Consume(ctx context.Context, in *sell.Goods) (*sell.GoodsInfo, error) {
	gi := s.GetGoods(in.Id)

	if gi == nil {
		return nil, fmt.Errorf("the goods %s not exist", in.Name)
	}
	if gi.Rest > 0 {
		gi.Rest--
		goods := &sell.GoodsInfo{
			G: &sell.Goods{
				Id:   gi.ID,
				Name: gi.Name,
			},
			Price: gi.Price,
			Rest:  gi.Rest,
		}
		fmt.Printf("%s sold, id= %d\n", in.Name, in.Id)
		return goods, nil
	} else {
		err := fmt.Errorf("%s empty rest", in.Name)
		return nil, err
	}
}

func (s *Shop) ListGoods(context.Context, *empty.Empty) (*sell.Totalgoods, error) {
	if s.goods == nil {
		return nil, fmt.Errorf("shop not init")
	}
	total := make([]*sell.GoodsInfo, 0)
	for _, good := range s.goods {
		total = append(total, &sell.GoodsInfo{
			G: &sell.Goods{
				Id:   good.ID,
				Name: good.Name,
			},
			Price: good.Price,
			Rest:  good.Rest,
		})
	}
	return &sell.Totalgoods{Total: total}, nil
}
