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


//
// Created by adi ENZEL on 8/21/19.
//

#ifndef E2_OPENTRACING_H
#define E2_OPENTRACING_H

#include <iostream>
#include <string>
#include <tracelibcpp/tracelibcpp.hpp>
#include <opentracing/tracer.h>
#include <opentracing/propagation.h>
#include <nlohmann/json.hpp>  // use nlohmann json library as an example

struct RICCarrierWriter : opentracing::TextMapWriter {
    explicit RICCarrierWriter(
            std::unordered_map<std::string, std::string>& data_)
            : data{data_} {}

    opentracing::expected<void> Set(
            opentracing::string_view key,
            opentracing::string_view value) const override {
        // OpenTracing uses opentracing::expected for error handling. This closely
        // follows the expected proposal for the C++ Standard Library. See
        //    http://open-std.org/JTC1/SC22/WG21/docs/papers/2017/p0323r3.pdf
        // for more background.
        opentracing::expected<void> result;

        auto was_successful = data.emplace(key, value);
        if (was_successful.second) {
            // Use a default constructed opentracing::expected<void> to indicate
            // success.
            return result;
        } else {
            // `key` clashes with existing data, so the span context can't be encoded
            // successfully; set opentracing::expected<void> to an std::error_code.
            return opentracing::make_unexpected(
                    std::make_error_code(std::errc::not_supported));
        }
    }

    std::unordered_map<std::string, std::string>& data;
};

struct RICCarrierReader : opentracing::TextMapReader {
    explicit RICCarrierReader(
            const std::unordered_map<std::string, std::string>& data_)
            : data{data_} {}

    using F = std::function<opentracing::expected<void>(
            opentracing::string_view, opentracing::string_view)>;

    opentracing::expected<void> ForeachKey(F f) const override {
        // Iterate through all key-value pairs, the tracer will use the relevant keys
        // to extract a span context.
        for (auto& key_value : data) {
            auto was_successful = f(key_value.first, key_value.second);
            if (!was_successful) {
                // If the callback returns and unexpected value, bail out of the loop.
                return was_successful;
            }
        }

        // Indicate successful iteration.
        return {};
    }

    // Optional, define TextMapReader::LookupKey to allow for faster extraction.
    opentracing::expected<opentracing::string_view> LookupKey(
            opentracing::string_view key) const override {
        auto iter = data.find(key);
        if (iter != data.end()) {
            return opentracing::make_unexpected(opentracing::key_not_found_error);
        }
        return opentracing::string_view{iter->second};
    }

    const std::unordered_map<std::string, std::string>& data;
};



#endif //E2_OPENTRACING_H
