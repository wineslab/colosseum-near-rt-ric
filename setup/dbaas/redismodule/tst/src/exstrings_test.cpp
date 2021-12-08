/*
 * Copyright (c) 2018-2020 Nokia.
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */


extern "C" {
#include "exstringsStub.h"
#include "redismodule.h"
}

#include "CppUTest/TestHarness.h"
#include "CppUTestExt/MockSupport.h"

#define OBJ_OP_NO 0
#define OBJ_OP_XX (1<<1)     /* OP if key exist */
#define OBJ_OP_NX (1<<2)     /* OP if key not exist */
#define OBJ_OP_IE (1<<4)     /* OP if equal old value */
#define OBJ_OP_NE (1<<5)     /* OP if not equal old value */

typedef struct RedisModuleBlockedClientArgs {
    RedisModuleBlockedClient *bc;
    RedisModuleString **argv;
    int argc;
} RedisModuleBlockedClientArgs;

TEST_GROUP(exstring)
{
    void setup()
    {

        mock().enable();
        mock().ignoreOtherCalls();
    }

    void teardown()
    {
        mock().clear();
        mock().disable();
    }

};

TEST(exstring, OnLoad)
{
    RedisModuleCtx ctx;
    int ret = RedisModule_OnLoad(&ctx, 0, 0);
    CHECK_EQUAL(ret, 0);
}

TEST(exstring, setie)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    int ret = SetIE_RedisCommand(&ctx, redisStrVec,  4);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);
    delete []redisStrVec;
}

TEST(exstring, setne)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    int ret = SetNE_RedisCommand(&ctx,redisStrVec, 4);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);
    delete []redisStrVec;
}

TEST(exstring, command_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    int ret = setStringGenericCommand(&ctx, 0, 3, OBJ_OP_IE);
    CHECK_EQUAL(ret, 1);
}


TEST(exstring, setie_command_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    mock().setData("RedisModule_OpenKey_no", 1);

    int ret = setStringGenericCommand(&ctx, redisStrVec, 4, OBJ_OP_IE);
    CHECK_EQUAL(ret, 0);
    delete []redisStrVec;

}


TEST(exstring, setie_command_has_key_set)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_set", 1);
    int ret = setStringGenericCommand(&ctx, redisStrVec, 4, OBJ_OP_IE);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 1);
    delete []redisStrVec;

}

TEST(exstring, setie_command_key_string_nosame)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);

    mock().setData("RedisModule_String_nosame", 1);


    int ret = setStringGenericCommand(&ctx, redisStrVec, 4, OBJ_OP_IE);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);
    delete []redisStrVec;

}

TEST(exstring, setie_command_key_same_string_reply)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    int ret = setStringGenericCommand(&ctx, redisStrVec, 4, OBJ_OP_IE);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    delete []redisStrVec;
}



TEST(exstring, setne_command_key_string_same_replrstr)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    //mock().setData("RedisModule_CallReplyType_str", 1);

    int ret = setStringGenericCommand(&ctx, redisStrVec, 4, OBJ_OP_NE);
    CHECK_EQUAL(ret, 0);
    delete []redisStrVec;

}

TEST(exstring, setne_command_setne_key_string_nosame_replrstr)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    int ret = setStringGenericCommand(&ctx, redisStrVec, 4, OBJ_OP_NE);
    CHECK_EQUAL(ret, 0);
    delete []redisStrVec;

}

TEST(exstring, delie)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIE_RedisCommand(&ctx, redisStrVec,  3);
    CHECK_EQUAL(ret, 0);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);
    delete []redisStrVec;
}

TEST(exstring, delne)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelNE_RedisCommand(&ctx,redisStrVec, 3);
    CHECK_EQUAL(ret, 0);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);
    delete []redisStrVec;
}

TEST(exstring, del_command_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    int ret = 0;
    ret = delStringGenericCommand(&ctx, 0, 2, OBJ_OP_IE);
    CHECK_EQUAL(ret, 1);

    ret = 0;
    ret = delStringGenericCommand(&ctx, 0, 4, OBJ_OP_NE);
    CHECK_EQUAL(ret, 1);
}

TEST(exstring, delie_command_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    mock().setData("RedisModule_OpenKey_no", 1);

    int ret = delStringGenericCommand(&ctx, redisStrVec, 3, OBJ_OP_IE);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithLongLong").getIntValue(), 0);
    delete []redisStrVec;

}

TEST(exstring, delie_command_have_key_set)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_set", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = delStringGenericCommand(&ctx, redisStrVec, 3, OBJ_OP_IE);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 1);
    delete []redisStrVec;
}

TEST(exstring, delie_command_key_string_nosame)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);

    mock().setData("RedisModule_String_nosame", 1);


    mock().expectOneCall("RedisModule_CloseKey");
    int ret = delStringGenericCommand(&ctx, redisStrVec, 3, OBJ_OP_IE);
    CHECK_EQUAL(ret, 0);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithLongLong").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);
    delete []redisStrVec;

}


TEST(exstring, delie_command_key_same_string_reply)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = delStringGenericCommand(&ctx, redisStrVec, 3, OBJ_OP_IE);
    CHECK_EQUAL(ret, 0);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);
    delete []redisStrVec;
}


TEST(exstring, delne_command_key_string_same_reply)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = delStringGenericCommand(&ctx, redisStrVec, 3, OBJ_OP_NE);
    CHECK_EQUAL(ret, 0);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithLongLong").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);
    delete []redisStrVec;

}

TEST(exstring, delne_command_key_string_nosame_reply)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[3]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = delStringGenericCommand(&ctx, redisStrVec, 3, OBJ_OP_NE);
    CHECK_EQUAL(ret, 0);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);
    delete []redisStrVec;

}

TEST(exstring, setpub_command_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    int ret = 0;

    ret = SetPub_RedisCommand(&ctx, 0, 2);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetPub_RedisCommand(&ctx, 0, 8);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetMPub_RedisCommand(&ctx, 0, 2);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetMPub_RedisCommand(&ctx, 0, 8);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetXXPub_RedisCommand(&ctx, 0, 3);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetXXPub_RedisCommand(&ctx, 0, 6);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetNXPub_RedisCommand(&ctx, 0, 3);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetNXMPub_RedisCommand(&ctx, 0, 3);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetNXPub_RedisCommand(&ctx, 0, 6);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetIEPub_RedisCommand(&ctx, 0, 4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetIEMPub_RedisCommand(&ctx, 0, 4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetIEPub_RedisCommand(&ctx, 0, 9);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetNEPub_RedisCommand(&ctx, 0, 4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = SetNEPub_RedisCommand(&ctx, 0, 9);
    CHECK_EQUAL(ret, REDISMODULE_ERR);
}

TEST(exstring, setpub_command_no_key_replystr)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    int ret = SetPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setmpub_command_negative_key_val_count)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[7]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;
    redisStrVec[6] = (RedisModuleString *)6;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_StringToLongLongCall_1", -1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = SetMPub_RedisCommand(&ctx, redisStrVec, 7);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(0, mock().getData("MSET").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setmpub_command_negative_chan_msg_count)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[7]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;
    redisStrVec[6] = (RedisModuleString *)6;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", -1);

    int ret = SetMPub_RedisCommand(&ctx, redisStrVec, 7);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(0, mock().getData("MSET").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setmpub_command_invalid_total_count)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[7]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;
    redisStrVec[6] = (RedisModuleString *)6;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_StringToLongLongCall_1", 100);
    mock().setData("RedisModule_StringToLongLongCall_2", 100);

    int ret = SetMPub_RedisCommand(&ctx, redisStrVec, 7);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(0, mock().getData("MSET").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setmpub_command_set)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[7]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;
    redisStrVec[6] = (RedisModuleString *)6;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = SetMPub_RedisCommand(&ctx, redisStrVec, 7);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(1, mock().getData("MSET").getIntValue());
    CHECK_EQUAL(1, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(2, mock().getData("RedisModule_FreeCallReply").getIntValue());
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setmpub_command_set_multipub)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[9]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;
    redisStrVec[6] = (RedisModuleString *)6;
    redisStrVec[7] = (RedisModuleString *)7;
    redisStrVec[8] = (RedisModuleString *)8;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 2);

    int ret = SetMPub_RedisCommand(&ctx, redisStrVec, 9);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(1, mock().getData("MSET").getIntValue());
    CHECK_EQUAL(2, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(3, mock().getData("RedisModule_FreeCallReply").getIntValue());
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setxxpub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetXXPub_RedisCommand(&ctx, redisStrVec, 5);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setxxpub_command_parameter_has_key_set)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;

    mock().setData("RedisModule_KeyType_set", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetXXPub_RedisCommand(&ctx, redisStrVec, 5);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setxxpub_command_has_key_string)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;

    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetXXPub_RedisCommand(&ctx, redisStrVec, 5);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}


TEST(exstring, setnxpub_command_has_key_string)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;

    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetNXPub_RedisCommand(&ctx, redisStrVec, 5);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setnxpub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetNXPub_RedisCommand(&ctx, redisStrVec, 5);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}



TEST(exstring, setiepub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetIEPub_RedisCommand(&ctx, redisStrVec, 6);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}


TEST(exstring, setiepub_command_key_string_nosame)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;

    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetIEPub_RedisCommand(&ctx, redisStrVec, 6);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setiepub_command_key_same_string_reply)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;

    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetIEPub_RedisCommand(&ctx, redisStrVec, 6);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 3);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setnepub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;

    mock().setData("RedisModule_KeyType_empty", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetNEPub_RedisCommand(&ctx, redisStrVec, 6);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 3);
    CHECK_EQUAL(mock().getData("RedisModule_AutoMemory").getIntValue(),1);

    delete []redisStrVec;

}

TEST(exstring, setnepub_command_key_string_same_reply)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;

    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetNEPub_RedisCommand(&ctx, redisStrVec, 6);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);

    delete []redisStrVec;

}


TEST(exstring, setnepub_command_key_string_nosame_reply)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)0;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)2;
    redisStrVec[3] = (RedisModuleString *)3;
    redisStrVec[4] = (RedisModuleString *)4;
    redisStrVec[5] = (RedisModuleString *)5;

    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = SetNEPub_RedisCommand(&ctx, redisStrVec, 6);

    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("MSET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 3);

    delete []redisStrVec;

}

TEST(exstring, delpub)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[4]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    int ret = DelPub_RedisCommand(&ctx, redisStrVec,  4);
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}

TEST(exstring, delmpub)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[8]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;
    redisStrVec[6] = (RedisModuleString *)1;
    redisStrVec[7] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 2);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec,  8);
    CHECK_EQUAL(ret, REDISMODULE_OK);

    delete []redisStrVec;
}

TEST(exstring, deliepub)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEPub_RedisCommand(&ctx, redisStrVec,  5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    delete []redisStrVec;
}

TEST(exstring, deliempub)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEMPub_RedisCommand(&ctx, redisStrVec,  5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    delete []redisStrVec;
}

TEST(exstring, delnepub)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelNEPub_RedisCommand(&ctx, redisStrVec,  5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    delete []redisStrVec;
}

TEST(exstring, delpub_command_parameter_number_incorrect)
{
    RedisModuleCtx ctx;
    int ret = 0;
    ret = DelPub_RedisCommand(&ctx, 0, 2);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = DelMPub_RedisCommand(&ctx, 0, 5);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = DelIEPub_RedisCommand(&ctx, 0, 4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = DelIEMPub_RedisCommand(&ctx, 0, 4);
    CHECK_EQUAL(ret, REDISMODULE_ERR);

    ret = 0;
    ret = DelNEPub_RedisCommand(&ctx, 0, 8);
    CHECK_EQUAL(ret, REDISMODULE_ERR);
}

TEST(exstring, delpub_command_reply_null)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    mock().setData("RedisModule_CallReplyInteger", 0);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_Call_Return_Null", 0);

    int ret = DelPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 1);
    delete []redisStrVec;

}

TEST(exstring, delpub_command_reply_error)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    mock().setData("RedisModule_CallReplyInteger", 0);
    mock().setData("RedisModule_CallReplyType_err", 1);

    int ret = DelPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    delete []redisStrVec;

}

TEST(exstring, delpub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    mock().setData("RedisModule_CallReplyInteger", 0);
    mock().setData("RedisModule_CallReplyType_inter", 1);

    int ret = DelPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_reply_null)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;

    mock().setData("RedisModule_Call_Return_Null", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(1, mock().getData("UNLINK").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_ReplyWithCallReply").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_reply_error)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 0);
    mock().setData("RedisModule_CallReplyType_err", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 0);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_key_deleted)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(1, mock().getData("UNLINK").getIntValue());
    CHECK_EQUAL(1, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithCallReply").getIntValue());
    CHECK_EQUAL(2, mock().getData("RedisModule_FreeCallReply").getIntValue());
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_key_deleted_multi_pub)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[10]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;
    redisStrVec[6] = (RedisModuleString *)1;
    redisStrVec[7] = (RedisModuleString *)1;
    redisStrVec[8] = (RedisModuleString *)1;
    redisStrVec[9] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", 3);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 10);
    CHECK_EQUAL(0, ret);
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(1, mock().getData("UNLINK").getIntValue());
    CHECK_EQUAL(3, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithCallReply").getIntValue());
    CHECK_EQUAL(4, mock().getData("RedisModule_FreeCallReply").getIntValue());
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_negative_del_count)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", -1);
    mock().setData("RedisModule_StringToLongLongCall_2", 1);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(0, ret);
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(0, mock().getData("UNLINK").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_negative_chan_msg_count)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 1);
    mock().setData("RedisModule_StringToLongLongCall_2", -1);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(0, ret);
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(0, mock().getData("UNLINK").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    delete []redisStrVec;

}

TEST(exstring, delmpub_command_invalid_total_count)
{
    RedisModuleCtx ctx;
    RedisModuleString **redisStrVec = new (RedisModuleString*[6]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    redisStrVec[5] = (RedisModuleString *)1;

    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().setData("RedisModule_CallReplyType_inter", 1);
    mock().setData("RedisModule_StringToLongLongCallCount", 0);
    mock().setData("RedisModule_StringToLongLongCall_1", 100);
    mock().setData("RedisModule_StringToLongLongCall_2", 100);

    int ret = DelMPub_RedisCommand(&ctx, redisStrVec, 6);
    CHECK_EQUAL(0, ret);
    CHECK_EQUAL(0, mock().getData("GET").getIntValue());
    CHECK_EQUAL(0, mock().getData("UNLINK").getIntValue());
    CHECK_EQUAL(0, mock().getData("PUBLISH").getIntValue());
    CHECK_EQUAL(1, mock().getData("RedisModule_ReplyWithError").getIntValue());
    CHECK_EQUAL(0, mock().getData("RedisModule_FreeCallReply").getIntValue());
    delete []redisStrVec;

}

TEST(exstring, deliepub_command_has_no_key)
{
    RedisModuleCtx ctx;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;
    mock().setData("RedisModule_KeyType_empty", 1);

    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithLongLong").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);
    delete []redisStrVec;

}

TEST(exstring, deliepub_command_has_key_set)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_set", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 0);

    delete []redisStrVec;
}

TEST(exstring, deliepub_command_key_string_nosame)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithLongLong").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);

    delete []redisStrVec;
}

TEST(exstring, deliepub_command_same_string_replynull)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_null", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithNull").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 2);

    delete []redisStrVec;
}

TEST(exstring, deliepub_command_same_string_reply)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelIEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 3);

    delete []redisStrVec;
}

TEST(exstring, delnepub_command_same_string_reply)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_same", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelNEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithLongLong").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 1);

    delete []redisStrVec;
}

TEST(exstring, delnepub_command_nosame_string_reply)
{
    RedisModuleCtx ctx;
    //RedisModuleString str;
    RedisModuleString ** redisStrVec = new (RedisModuleString*[5]);

    redisStrVec[0] = (RedisModuleString *)1;
    redisStrVec[1] = (RedisModuleString *)1;
    redisStrVec[2] = (RedisModuleString *)1;
    redisStrVec[3] = (RedisModuleString *)1;
    redisStrVec[4] = (RedisModuleString *)1;

    mock().setData("RedisModule_OpenKey_have", 1);
    mock().setData("RedisModule_KeyType_str", 1);
    mock().setData("RedisModule_String_nosame", 1);
    mock().setData("RedisModule_CallReplyType_str", 1);
    mock().setData("RedisModule_CallReplyInteger", 1);
    mock().expectOneCall("RedisModule_CloseKey");
    int ret = DelNEPub_RedisCommand(&ctx, redisStrVec, 5);
    CHECK_EQUAL(ret, REDISMODULE_OK);
    mock().checkExpectations();
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithError").getIntValue(), 0);
    CHECK_EQUAL(mock().getData("GET").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("UNLINK").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("PUBLISH").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_ReplyWithCallReply").getIntValue(), 1);
    CHECK_EQUAL(mock().getData("RedisModule_FreeCallReply").getIntValue(), 3);

    delete []redisStrVec;
}
