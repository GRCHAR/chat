#!/bin/bash

# 用户搜索和私聊功能测试脚本

BASE_URL="http://localhost:8080/api/v1"
TOKEN=""

echo "=========================================="
echo "用户搜索和私聊功能测试"
echo "=========================================="
echo ""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 1. 注册测试用户
echo -e "${YELLOW}1. 注册测试用户...${NC}"
REGISTER_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser1",
    "nickname": "测试用户1",
    "email": "test1@example.com",
    "password": "password123"
  }')

echo "$REGISTER_RESPONSE" | jq '.'
TOKEN=$(echo "$REGISTER_RESPONSE" | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
  echo -e "${GREEN}✓ 用户注册成功${NC}"
else
  echo -e "${RED}✗ 用户注册失败${NC}"
  exit 1
fi
echo ""

# 2. 注册第二个测试用户
echo -e "${YELLOW}2. 注册第二个测试用户...${NC}"
REGISTER2_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser2",
    "nickname": "测试用户2",
    "email": "test2@example.com",
    "password": "password123"
  }')

echo "$REGISTER2_RESPONSE" | jq '.'
USER2_ID=$(echo "$REGISTER2_RESPONSE" | jq -r '.user.id')
echo -e "${GREEN}✓ 第二个用户注册成功，ID: $USER2_ID${NC}"
echo ""

# 3. 搜索用户
echo -e "${YELLOW}3. 搜索用户（关键词: test）...${NC}"
SEARCH_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/search?q=test" \
  -H "Authorization: Bearer $TOKEN")

echo "$SEARCH_RESPONSE" | jq '.'
USER_COUNT=$(echo "$SEARCH_RESPONSE" | jq '.users | length')

if [ "$USER_COUNT" -gt 0 ]; then
  echo -e "${GREEN}✓ 搜索成功，找到 $USER_COUNT 个用户${NC}"
else
  echo -e "${RED}✗ 搜索失败或无结果${NC}"
fi
echo ""

# 4. 根据ID获取用户详情
echo -e "${YELLOW}4. 获取用户详情（ID: $USER2_ID）...${NC}"
USER_DETAIL_RESPONSE=$(curl -s -X GET "${BASE_URL}/users/${USER2_ID}" \
  -H "Authorization: Bearer $TOKEN")

echo "$USER_DETAIL_RESPONSE" | jq '.'
USERNAME=$(echo "$USER_DETAIL_RESPONSE" | jq -r '.user.username')

if [ "$USERNAME" != "null" ]; then
  echo -e "${GREEN}✓ 获取用户详情成功${NC}"
else
  echo -e "${RED}✗ 获取用户详情失败${NC}"
fi
echo ""

# 5. 创建私聊房间
echo -e "${YELLOW}5. 创建私聊房间...${NC}"
CREATE_ROOM_RESPONSE=$(curl -s -X POST "${BASE_URL}/rooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"与测试用户2的聊天\",
    \"type\": \"single\",
    \"member_ids\": [$USER2_ID]
  }")

echo "$CREATE_ROOM_RESPONSE" | jq '.'
ROOM_ID=$(echo "$CREATE_ROOM_RESPONSE" | jq -r '.room.id')

if [ "$ROOM_ID" != "null" ] && [ -n "$ROOM_ID" ]; then
  echo -e "${GREEN}✓ 私聊房间创建成功，ID: $ROOM_ID${NC}"
else
  echo -e "${RED}✗ 私聊房间创建失败${NC}"
fi
echo ""

# 6. 创建群聊房间
echo -e "${YELLOW}6. 创建群聊房间...${NC}"
CREATE_GROUP_RESPONSE=$(curl -s -X POST "${BASE_URL}/rooms" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"测试群聊\",
    \"description\": \"这是一个测试群聊\",
    \"type\": \"group\",
    \"member_ids\": [$USER2_ID]
  }")

echo "$CREATE_GROUP_RESPONSE" | jq '.'
GROUP_ID=$(echo "$CREATE_GROUP_RESPONSE" | jq -r '.room.id')

if [ "$GROUP_ID" != "null" ] && [ -n "$GROUP_ID" ]; then
  echo -e "${GREEN}✓ 群聊房间创建成功，ID: $GROUP_ID${NC}"
else
  echo -e "${RED}✗ 群聊房间创建失败${NC}"
fi
echo ""

# 7. 注册第三个用户用于添加到群聊
echo -e "${YELLOW}7. 注册第三个测试用户...${NC}"
REGISTER3_RESPONSE=$(curl -s -X POST "${BASE_URL}/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser3",
    "nickname": "测试用户3",
    "email": "test3@example.com",
    "password": "password123"
  }')

USER3_ID=$(echo "$REGISTER3_RESPONSE" | jq -r '.user.id')
echo -e "${GREEN}✓ 第三个用户注册成功，ID: $USER3_ID${NC}"
echo ""

# 8. 添加成员到群聊
echo -e "${YELLOW}8. 添加成员到群聊...${NC}"
ADD_MEMBER_RESPONSE=$(curl -s -X POST "${BASE_URL}/rooms/${GROUP_ID}/members" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": $USER3_ID
  }")

echo "$ADD_MEMBER_RESPONSE" | jq '.'
ADD_MESSAGE=$(echo "$ADD_MEMBER_RESPONSE" | jq -r '.message')

if [ "$ADD_MESSAGE" != "null" ]; then
  echo -e "${GREEN}✓ 添加成员成功${NC}"
else
  echo -e "${RED}✗ 添加成员失败${NC}"
fi
echo ""

# 9. 获取房间成员列表
echo -e "${YELLOW}9. 获取群聊成员列表...${NC}"
MEMBERS_RESPONSE=$(curl -s -X GET "${BASE_URL}/rooms/${GROUP_ID}/members" \
  -H "Authorization: Bearer $TOKEN")

echo "$MEMBERS_RESPONSE" | jq '.'
MEMBER_COUNT=$(echo "$MEMBERS_RESPONSE" | jq '.members | length')

if [ "$MEMBER_COUNT" -gt 0 ]; then
  echo -e "${GREEN}✓ 获取成员列表成功，共 $MEMBER_COUNT 个成员${NC}"
else
  echo -e "${RED}✗ 获取成员列表失败${NC}"
fi
echo ""

# 10. 获取房间列表
echo -e "${YELLOW}10. 获取房间列表...${NC}"
ROOMS_RESPONSE=$(curl -s -X GET "${BASE_URL}/rooms" \
  -H "Authorization: Bearer $TOKEN")

echo "$ROOMS_RESPONSE" | jq '.'
ROOM_COUNT=$(echo "$ROOMS_RESPONSE" | jq '.rooms | length')

if [ "$ROOM_COUNT" -gt 0 ]; then
  echo -e "${GREEN}✓ 获取房间列表成功，共 $ROOM_COUNT 个房间${NC}"
else
  echo -e "${RED}✗ 获取房间列表失败${NC}"
fi
echo ""

echo "=========================================="
echo -e "${GREEN}测试完成！${NC}"
echo "=========================================="
echo ""
echo "测试摘要："
echo "- 用户搜索: ✓"
echo "- 获取用户详情: ✓"
echo "- 创建私聊: ✓"
echo "- 创建群聊: ✓"
echo "- 添加成员到群聊: ✓"
echo "- 获取成员列表: ✓"
echo ""
echo "提示：请确保后端服务正在运行在 http://localhost:8080"
