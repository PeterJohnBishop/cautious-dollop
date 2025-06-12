//
//  Item.swift
//  swift-cautious-dollop
//
//  Created by TheComputer on 6/12/25.
//

import Foundation
import SwiftData

@Model
final class Item {
    var timestamp: Date
    
    init(timestamp: Date) {
        self.timestamp = timestamp
    }
}
