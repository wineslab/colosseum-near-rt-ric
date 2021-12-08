/*
     * Copyright 2019 AT&T Intellectual Property
     * Copyright 2019 Nokia
     *
     * Licensed under the Apache License, Version 2.0 (the "License");
     * you may not use this file except in compliance with the License.
     * You may obtain a copy of the License at
     *
     *      http://www.apache.org/licenses/LICENSE-2.0
     *
     * Unless required by applicable law or agreed to in writing, software
     * distributed under the License is distributed on an "AS IS" BASIS,
     * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
     * See the License for the specific language governing permissions and
     * limitations under the License.
 */

/*
 * This source code is part of the near-RT RIC (RAN Intelligent Controller)
 * platform project (RICP).
 */



    #include <stdbool.h>
    #include <stdio.h>
    #include <stdlib.h>
    #include <x2reset_response_wrapper.h>

    void test_build_pack_x2reset_response();
    void test_unpack(void);

    int
    main(int argc, char* argv[])
    {
        test_build_pack_x2reset_response();
        exit(0);
    }

    void test_build_pack_x2reset_response(){
        size_t error_buf_size = 8192;
        size_t packed_buf_size = 4096;
        unsigned char responseDataBuf[packed_buf_size];
        char responseErrorBuf[error_buf_size];
        bool result = build_pack_x2reset_response(&packed_buf_size, responseDataBuf, error_buf_size, responseErrorBuf);

        if (!result) {
            printf("#test_build_pack_x2reset_response failed. Packing error %s\n", responseErrorBuf);
            return;
        }
        printf("x2reset response packed size:%lu\nPayload:\n", packed_buf_size);
        for (size_t i = 0; i < packed_buf_size; ++i)
            printf("%02x",responseDataBuf[i]);
        printf("\n");
    }

