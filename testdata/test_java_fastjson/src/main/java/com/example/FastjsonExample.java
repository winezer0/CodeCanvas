package com.example;

import com.alibaba.fastjson.JSON;
import com.alibaba.fastjson.JSONObject;

/**
 * Fastjson 使用示例
 */
public class FastjsonExample {
    
    public static void main(String[] args) {
        // 创建一个Java对象
        User user = new User();
        user.setId(1);
        user.setName("测试用户");
        user.setAge(30);
        
        // 使用fastjson将Java对象转换为JSON字符串
        String jsonString = JSON.toJSONString(user);
        System.out.println("JSON字符串: " + jsonString);
        
        // 使用fastjson将JSON字符串转换为Java对象
        User parsedUser = JSON.parseObject(jsonString, User.class);
        System.out.println("解析后的用户: " + parsedUser.getName());
        
        // 直接使用JSONObject
        JSONObject jsonObject = new JSONObject();
        jsonObject.put("key", "value");
        System.out.println("JSONObject: " + jsonObject.toJSONString());
    }
    
    /**
     * 用户实体类
     */
    static class User {
        private int id;
        private String name;
        private int age;
        
        // getter和setter方法
        public int getId() {
            return id;
        }
        
        public void setId(int id) {
            this.id = id;
        }
        
        public String getName() {
            return name;
        }
        
        public void setName(String name) {
            this.name = name;
        }
        
        public int getAge() {
            return age;
        }
        
        public void setAge(int age) {
            this.age = age;
        }
    }
}