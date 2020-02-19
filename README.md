#FYP-BACKEND
 - backed service of the fyp system
 ###1.MODELS
 ####ORDER: order_id=int64, table_id=int8, user_id=int64, items=[]Item, total_price=float32, is_payed=bool, is_ready=bool
 ####ITEM: item_id=int64, item_name=string, item_description=string, item_img=string, item_price=float32
 ####USER: user_id=int64, email=string, password=string, orders=[]Order
 
 - protoc --proto_path=proto -I/--proto_path=proto --go_out=plugins=grpc:models user.proto order.proto item.proto
 
 staff members as users will only be crated from the managers side, the same applies to the items  
 
 